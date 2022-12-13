package cookie

import (
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type UserSessionStore struct {
	codecs     []securecookie.Codec
	options    *sessions.Options
	repo       repository.UserSessionRepository
	cookieName string
}

type UserSessionStoreOpts struct {
	SessionRepository   repository.UserSessionRepository
	CookieSecrets       []string
	CookieAllowInsecure bool
	CookieDomain        string
	CookieName          string
}

func NewUserSessionStore(opts *UserSessionStoreOpts) (*UserSessionStore, error) {
	keyPairs := [][]byte{}

	for _, key := range opts.CookieSecrets {
		keyPairs = append(keyPairs, []byte(key))
	}

	res := &UserSessionStore{
		codecs: securecookie.CodecsFromPairs(keyPairs...),
		options: &sessions.Options{
			Path:     "/",
			Domain:   opts.CookieDomain,
			MaxAge:   86400 * 30,
			Secure:   !opts.CookieAllowInsecure,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
		repo:       opts.SessionRepository,
		cookieName: opts.CookieName,
	}

	return res, nil
}

func (store *UserSessionStore) GetName() string {
	return store.cookieName
}

func (store *UserSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(store, name)

	if session == nil {
		return nil, nil
	}

	opts := *store.options
	session.Options = &(opts)
	session.IsNew = true

	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, store.codecs...)
		if err == nil {
			err = store.load(session)

			if err != nil {
				if errors.Is(err, repository.RepositoryErrorNotFound) {
					err = nil
				} else if strings.Contains(err.Error(), "expired timestamp") {
					err = nil
					session.IsNew = false
				}
			} else {
				session.IsNew = false
			}
		}
	}

	store.MaxAge(store.options.MaxAge)

	return session, err
}

// Get Fetches a session for a given name after it has been added to the
// registry.
func (store *UserSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(store, name)
}

// Save saves the given session into the database and deletes cookies if needed
func (store *UserSessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	repo := store.repo

	// Set delete if max-age is < 0
	if session.Options.MaxAge < 0 {
		if _, err := repo.DeleteUserSession(&models.UserSession{Key: session.ID}); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Generate a random session ID key suitable for storage in the DB
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32),
			), "=")
	}

	if err := store.save(session); err != nil {
		return err
	}

	// Keep the session ID key in a cookie so it can be looked up in DB later.
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, store.codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// save writes encoded session.Values to a database record.
// writes to http_sessions table by default.
func (store *UserSessionStore) save(session *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values, store.codecs...)
	if err != nil {
		return err
	}

	exOn := session.Values["expires_on"]

	var expiresOn time.Time

	if exOn == nil {
		expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	} else {
		expiresOn = exOn.(time.Time)
		if expiresOn.Sub(time.Now().Add(time.Second*time.Duration(session.Options.MaxAge))) < 0 {
			expiresOn = time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
		}
	}

	s := &models.UserSession{
		Key:       session.ID,
		Data:      []byte(encoded),
		ExpiresAt: expiresOn,
	}

	repo := store.repo

	if session.IsNew {
		_, createErr := repo.CreateUserSession(s)
		return createErr
	}

	_, updateErr := repo.UpdateUserSession(s)
	return updateErr
}

// load fetches a session by ID from the database and decodes its content
// into session.Values.
func (store *UserSessionStore) load(session *sessions.Session) error {
	res, err := store.repo.ReadUserSessionByKey(session.ID)

	if err != nil {
		return err
	}

	return securecookie.DecodeMulti(session.Name(), string(res.Data), &session.Values, store.codecs...)
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new PGStore is 4096. PostgreSQL allows for max
// value sizes of up to 1GB (http://www.postgresql.org/docs/current/interactive/datatype-character.html)
func (store *UserSessionStore) MaxLength(l int) {
	for _, c := range store.codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (store *UserSessionStore) MaxAge(age int) {
	store.options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range store.codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

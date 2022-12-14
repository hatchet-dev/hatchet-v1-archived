// go:build test
package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/hatchet-dev/hatchet/internal/adapter"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	"github.com/hatchet-dev/hatchet/internal/migrate"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm"
)

type Tester struct {
	key        *[32]byte
	dbFileName string
	conf       *database.Config
	initData   *InitData
}

func NewTestEnv(t *testing.T) *Tester {
	t.Helper()

	tester := new(Tester)

	// generate a random string for the db file name
	randName, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.dbFileName = fmt.Sprintf("./%s.db", randName)

	// generate random bytes for the encryption key
	key := encryption.NewEncryptionKey()

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	db, err := adapter.New(&database.ConfigFile{
		EncryptionKey: string(key[:]),
		SQLLite:       true,
		SQLLitePath:   tester.dbFileName,
	})

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	// call automigration
	err = migrate.AutoMigrate(db, false)

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	tester.key = key

	tester.conf = &database.Config{
		GormDB:     db,
		Repository: gorm.NewRepository(db, key),
	}

	tester.initData = new(InitData)

	return tester
}

func Cleanup(t *testing.T, tester *Tester) {
	t.Helper()

	// remove the created file file
	os.Remove(tester.dbFileName)
}

func RunTestWithDatabase(t *testing.T, test func(config *database.Config) error, initMethods ...InitDataFunc) {
	t.Helper()

	tester := NewTestEnv(t)
	defer Cleanup(t, tester)

	for _, f := range initMethods {
		err := f(t, tester.conf)

		if err != nil {
			t.Fatalf("could not init data: %v\n", err)
		}
	}

	err := test(tester.conf)

	if err != nil {
		t.Fatalf("%v\n", err)
	}
}

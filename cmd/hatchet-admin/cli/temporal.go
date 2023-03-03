package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"

	"github.com/spf13/cobra"

	hatchettoken "github.com/hatchet-dev/hatchet/internal/auth/token"
)

var temporalCmd = &cobra.Command{
	Use:   "temporal",
	Short: "Command used to manage temporal settings and tokens",
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create-token",
	Short: "Creates a new internal Temporal token",
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreateTemporalToken()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [temporal create-token]:", err.Error())
			os.Exit(1)
		}
	},
}

var tokenCreateSigningKey string
var tokenIssuer string = "http://127.0.0.1:7223"
var tokenAudience string = "http://127.0.0.1:7223"

func init() {
	rootCmd.AddCommand(temporalCmd)
	temporalCmd.AddCommand(tokenCreateCmd)

	tokenCreateCmd.PersistentFlags().StringVar(
		&tokenCreateSigningKey,
		"signing-key",
		"",
		"The signing key for the token.",
	)

	tokenCreateCmd.MarkPersistentFlagRequired("signing-key")

	tokenCreateCmd.PersistentFlags().StringVar(
		&tokenIssuer,
		"iss",
		"http://127.0.0.1:7223",
		"The issuer URL for the token (includes the protocol).",
	)

	tokenCreateCmd.PersistentFlags().StringVar(
		&tokenAudience,
		"aud",
		"http://127.0.0.1:7223",
		"The audience URL for the token (includes the protocol).",
	)
}

func runCreateTemporalToken() error {
	internalAuthConfig := &temporal.InternalAuthConfig{
		InternalNamespace:  temporal.DefaultInternalNamespace,
		InternalSigningKey: []byte(tokenCreateSigningKey),
		InternalTokenOpts: hatchettoken.TokenOpts{
			Issuer:   tokenIssuer,
			Audience: []string{tokenAudience},
		},
	}

	token, err := token.GenerateInternalToken(internalAuthConfig)

	if err != nil {
		return err
	}

	fmt.Println(token)

	return nil
}

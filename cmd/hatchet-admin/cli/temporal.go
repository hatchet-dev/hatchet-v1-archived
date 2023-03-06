package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"

	"github.com/spf13/cobra"

	hatchettoken "github.com/hatchet-dev/hatchet/internal/auth/token"
)

var temporalCmd = &cobra.Command{
	Use:   "temporal",
	Short: "Command used to manage temporal settings and tokens",
}

var internalTokenCreateCmd = &cobra.Command{
	Use:   "create-internal-token",
	Short: "Creates a new internal Temporal token",
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreateTemporalInternalToken()

		if err != nil {
			red := color.New(color.FgRed)
			red.Printf("Error running [%s]:%s\n", cmd.Use, err.Error())
			os.Exit(1)
		}
	},
}

var workerTokenCreateCmd = &cobra.Command{
	Use:   "create-worker-token",
	Short: "Creates a new worker Temporal token",
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreateTemporalWorkerToken()

		if err != nil {
			red := color.New(color.FgRed)
			red.Printf("Error running [%s]:%s\n", cmd.Use, err.Error())
			os.Exit(1)
		}
	},
}

var tokenCreateSigningKey string
var tokenIssuer string
var tokenAudience string

var workerTokenTeamID string

func init() {
	rootCmd.AddCommand(temporalCmd)
	temporalCmd.AddCommand(internalTokenCreateCmd)
	temporalCmd.AddCommand(workerTokenCreateCmd)

	internalTokenCreateCmd.PersistentFlags().StringVar(
		&tokenCreateSigningKey,
		"signing-key",
		"",
		"The signing key for the token.",
	)

	internalTokenCreateCmd.MarkPersistentFlagRequired("signing-key")

	internalTokenCreateCmd.PersistentFlags().StringVar(
		&tokenIssuer,
		"iss",
		sc.Config.TemporalClient.GetBroadcastAddress(),
		"The issuer URL for the token (includes the protocol).",
	)

	internalTokenCreateCmd.PersistentFlags().StringVar(
		&tokenAudience,
		"aud",
		sc.Config.TemporalClient.GetBroadcastAddress(),
		"The audience URL for the token (includes the protocol).",
	)

	workerTokenCreateCmd.PersistentFlags().StringVar(
		&workerTokenTeamID,
		"team-id",
		"",
		"The team id to generate the worker token for",
	)

	workerTokenCreateCmd.MarkPersistentFlagRequired("team-id")
}

func runCreateTemporalInternalToken() error {
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

func runCreateTemporalWorkerToken() error {
	wt, err := models.NewWorkerTokenFromTeamID(workerTokenTeamID)

	if err != nil {
		return err
	}

	token, err := hatchettoken.GenerateTokenFromWT(wt, &hatchettoken.TokenOpts{
		Issuer:   tokenIssuer,
		Audience: []string{tokenAudience},
	})

	if err != nil {
		return err
	}

	wt, err = sc.DB.Repository.WorkerToken().CreateWorkerToken(wt)

	if err != nil {
		return err
	}

	fmt.Println(token)

	return nil
}

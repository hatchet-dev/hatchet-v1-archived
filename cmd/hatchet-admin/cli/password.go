package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/internal/models"

	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Command used to manage Hatchet instance passwords",
}

var passwordResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Creates a new default monitor for a given team",
	Run: func(cmd *cobra.Command, args []string) {
		err := runResetPassword()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [password reset]:", err.Error())
			os.Exit(1)
		}
	},
}

type passwordResetRequest struct {
	Email string `form:"required,max=255,email"`
}

var resetRequest = passwordResetRequest{}

func init() {
	rootCmd.AddCommand(passwordCmd)
	passwordCmd.AddCommand(passwordResetCmd)

	passwordResetCmd.PersistentFlags().StringVar(
		&resetRequest.Email,
		"email",
		"",
		"The email for the user.",
	)

	passwordResetCmd.MarkPersistentFlagRequired("email")
}

func runResetPassword() error {
	// validate the monitor
	v := handlerutils.NewDefaultValidator()

	reqErr := v.Validate(resetRequest)

	if reqErr != nil {
		return reqErr
	}

	_, err := sc.DB.Repository.User().ReadUserByEmail(resetRequest.Email)

	if err != nil {
		return err
	}

	pwResetToken, err := models.NewPasswordResetTokenFromEmail(resetRequest.Email)

	if err != nil {
		return err
	}

	// this is the only time we'll recover the raw pw reset token, so we store it in the values
	queryVals := url.Values{
		"token": []string{string(pwResetToken.Token)},
		"email": []string{resetRequest.Email},
	}

	pwResetToken, err = sc.DB.Repository.PasswordResetToken().CreatePasswordResetToken(pwResetToken)

	if err != nil {
		return err
	}

	queryVals.Add("token_id", pwResetToken.ID)

	fmt.Printf("Share this password reset link with your team member: %s/reset_password/finalize?%s\n", sc.ServerRuntimeConfig.ServerURL, queryVals.Encode())

	return nil
}

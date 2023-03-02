package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/internal/models"

	"github.com/spf13/cobra"
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Command used to manage notification settings for the Hatchet instance",
}

var inboxCreateCmd = &cobra.Command{
	Use:   "create-inbox",
	Short: "Creates a new inbox for a given team",
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreateInbox()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor create-default]:", err.Error())
			os.Exit(1)
		}
	},
}

type inboxCreateRequest struct {
	TeamID string `form:"required,uuid"`
}

var inboxCreateReq = inboxCreateRequest{}

func init() {
	rootCmd.AddCommand(notificationCmd)
	notificationCmd.AddCommand(inboxCreateCmd)

	inboxCreateCmd.PersistentFlags().StringVar(
		&inboxCreateReq.TeamID,
		"team",
		"",
		"The team ID to assign this inbox to.",
	)

	inboxCreateCmd.MarkPersistentFlagRequired("team")
}

func runCreateInbox() error {
	// validate the monitor
	v := handlerutils.NewDefaultValidator()

	reqErr := v.Validate(inboxCreateReq)

	if reqErr != nil {
		return reqErr
	}

	inbox := &models.NotificationInbox{
		TeamID: inboxCreateReq.TeamID,
	}

	inbox, err := sc.DB.Repository.Notification().CreateNotificationInbox(inbox)

	if err != nil {
		return err
	}

	return err
}

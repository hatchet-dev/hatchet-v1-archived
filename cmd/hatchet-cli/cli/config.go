package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage Hatchet CLI configuration options",
}

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "displays the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		err := printConfig()

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not get config: %v\n", err)
			os.Exit(1)
		}
	},
}

var setAddressCmd = &cobra.Command{
	Use:   "set-address",
	Short: "set the address of the Hatchet instance",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setAddress(args[0])

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not set address: %v\n", err)
			os.Exit(1)
		}
	},
}

var setOrganizationCmd = &cobra.Command{
	Use:   "set-organization",
	Short: "set the default organization of the Hatchet instance",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setOrganization(args[0])

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not set organization: %v\n", err)
			os.Exit(1)
		}
	},
}

var setTeamCmd = &cobra.Command{
	Use:   "set-team",
	Short: "set the default team of the Hatchet instance",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setTeam(args[0])

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not set team: %v\n", err)
			os.Exit(1)
		}
	},
}

var setAPITokenCmd = &cobra.Command{
	Use:   "set-api-token",
	Short: "set the default api token of the Hatchet instance",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := setAPIToken(args[0])

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not set api token: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(setAddressCmd)
	configCmd.AddCommand(setOrganizationCmd)
	configCmd.AddCommand(setTeamCmd)
	configCmd.AddCommand(setAPITokenCmd)
}

func setAddress(address string) error {
	address = strings.TrimRight(address, "/")

	v.Set("address", address)
	color.New(color.FgGreen).Printf("Set the current address as %s\n", address)
	err := v.WriteConfig()

	if err != nil {
		return err
	}

	return nil
}

func setOrganization(orgID string) error {
	v.Set("organizationID", orgID)
	color.New(color.FgGreen).Printf("Set the current organization as %s\n", orgID)
	err := v.WriteConfig()

	if err != nil {
		return err
	}

	return nil
}

func setTeam(teamID string) error {
	v.Set("teamID", teamID)
	color.New(color.FgGreen).Printf("Set the current team as %s\n", teamID)
	err := v.WriteConfig()

	if err != nil {
		return err
	}

	return nil
}

func setAPIToken(apiToken string) error {
	v.Set("apiToken", apiToken)
	color.New(color.FgGreen).Printf("Set the current api token successfully\n")
	err := v.WriteConfig()

	if err != nil {
		return err
	}

	if apiToken == "" {
		return nil
	}

	// reload api client with new token
	clientConf := swagger.NewConfiguration()

	clientConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	config.APIClient = swagger.NewAPIClient(clientConf)

	// list user organizations to determine if org id needs to be reset
	orgs, _, err := config.APIClient.UsersApi.ListUserOrganizations(
		context.Background(),
		&swagger.UsersApiListUserOrganizationsOpts{},
	)

	if err != nil {
		return fmt.Errorf("could not list user organizations with this token: %w", err)
	}

	currOrgID := v.GetString("organizationID")

	var shouldSetOrgID = currOrgID == ""

	if currOrgID != "" {
		foundOrgID := false

		for _, org := range orgs.Rows {
			if org.Id == currOrgID {
				foundOrgID = true
				break
			}
		}

		if !foundOrgID {
			shouldSetOrgID = true
		}
	}

	if shouldSetOrgID {
		if len(orgs.Rows) > 0 {
			err = setOrganization(orgs.Rows[0].Id)

			if err != nil {
				return err
			}
		}
	}

	currTeamID := v.GetString("teamID")

	shouldSetTeamID := currTeamID == ""

	teams, _, err := config.APIClient.UsersApi.ListUserTeams(
		context.Background(),
		&swagger.UsersApiListUserTeamsOpts{},
	)

	if err != nil {
		return fmt.Errorf("could not list user teams with this token: %w", err)
	}

	if currTeamID != "" {
		foundTeamID := false

		for _, team := range teams.Rows {
			if team.Id == currTeamID {
				foundTeamID = true
				break
			}
		}

		if !foundTeamID {
			shouldSetTeamID = true
		}
	}

	if shouldSetTeamID {
		if len(teams.Rows) > 0 {
			err = setTeam(teams.Rows[0].Id)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func printConfig() error {
	config, err := ioutil.ReadFile(filepath.Join(home, ".hatchet", "hatchet.yaml"))

	if err != nil {
		return err
	}

	fmt.Print(string(config))

	return nil
}

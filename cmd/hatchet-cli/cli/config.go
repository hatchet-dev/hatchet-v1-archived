package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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

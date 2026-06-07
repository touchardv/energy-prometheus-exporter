package cmd

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/touchardv/energy-prometheus-exporter/pkg/homewizard"
)

const (
	flagHomewizardURL      = "homewizard-url"
	flagHomewizardToken    = "homewizard-token"
	flagHomewizardUsername = "homewizard-username"
)

func newHomewizardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "homewizard",
		Short: "Manage a Homewizard device",
	}
	cmd.PersistentFlags().StringP(flagHomewizardURL, "", "", "The URL of the Homewizard device")

	createLocalUser := &cobra.Command{
		Use: "create-user",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Parent().Flags().GetString(flagHomewizardURL)
			if err != nil {
				return err
			}
			username, _ := cmd.Flags().GetString(flagHomewizardUsername)
			client := homewizard.NewAPIv2Client(url)
			return client.CreateLocalUser(username)
		},
	}
	createLocalUser.Flags().StringP(flagHomewizardUsername, "", "", "The name of the user to register in the Homewizard device")

	listUsers := &cobra.Command{
		Use:   "list-users",
		Short: "List registered user(s)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Parent().Flags().GetString(flagHomewizardURL)
			if err != nil {
				return err
			}
			token, err := cmd.Flags().GetString(flagHomewizardToken)
			if err != nil {
				return err
			}
			client := homewizard.NewAuthenticatedAPIv2Client(url, token)
			users, err := client.ListUsers()
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "% 20s %s\n", "NAME", "CURRENT")
			for _, user := range users {
				fmt.Fprintf(os.Stdout, "% 20s %t\n", user.Name, user.Current)
			}
			return nil
		},
	}
	listUsers.Flags().StringP(flagHomewizardToken, "", "", "Homewizard user authentication token")

	cmd.AddCommand(createLocalUser)
	cmd.AddCommand(listUsers)
	return cmd
}

func addHomewizardCollectorFlags(fs *pflag.FlagSet) {
	fs.StringP(flagHomewizardURL, "", "", "Homewizard device URL")
	fs.StringP(flagHomewizardToken, "", "", "Homewizard user authentication token")
}

func addHomewizardCollector(fs *pflag.FlagSet, reg *prometheus.Registry) error {
	url, err := fs.GetString(flagHomewizardURL)
	if err != nil {
		return err
	}
	if url == "" {
		log.Println("homewizard: url not set, collector disabled")
		return nil
	}
	log.Printf("homewizard: collector enabled (url: %s)", url)
	token, err := fs.GetString(flagHomewizardToken)
	if err != nil {
		return err
	}
	client := homewizard.NewAuthenticatedAPIv2Client(url, token)
	reg.MustRegister(client)
	return nil
}

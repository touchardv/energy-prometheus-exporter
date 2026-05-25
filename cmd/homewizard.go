package cmd

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/touchardv/energy-prometheus-exporter/pkg/homewizard"
)

const (
	prefix           = "homewizard"
	urlProperty      = "url"
	tokenProperty    = "token"
	usernameProperty = "username"
)

func registerHomewizard(cmd *cobra.Command) {
	cmd.Flags().StringP(propertyName(urlProperty), "", "", "The URL of the Homewizard device")
	cmd.Flags().StringP(propertyName(tokenProperty), "", "", "The user authentication token")
}

func buildCollectorForHomewizard(cmd *cobra.Command) prometheus.Collector {
	url, err := cmd.Flags().GetString(propertyName(urlProperty))
	if err != nil {
		return nil
	}
	token, err := cmd.Flags().GetString(propertyName(tokenProperty))
	if err != nil {
		return nil
	}
	return homewizard.NewAPIv2Client(url, homewizard.WithToken(token))
}

func propertyName(n string) string {
	return fmt.Sprintf("%s-%s", prefix, n)
}

func newHomewizardCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "homewizard",
		Short: "Manage a Homewizard device",
	}
	cmd.PersistentFlags().StringP(urlProperty, "", "", "The URL of the Homewizard device")

	createLocalUser := &cobra.Command{
		Use: "create-user",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd, prefix)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Parent().Flags().GetString(urlProperty)
			if err != nil {
				return err
			}
			username, _ := cmd.Flags().GetString(usernameProperty)
			client := homewizard.NewAPIv2Client(url)
			return client.CreateLocalUser(username)
		},
	}
	createLocalUser.Flags().StringP(usernameProperty, "", "", "The name of the user to register in the Homewizard device")

	listUsers := &cobra.Command{
		Use:   "list-users",
		Short: "List registered user(s)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd, prefix)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Parent().Flags().GetString(urlProperty)
			if err != nil {
				return err
			}
			token, err := cmd.Flags().GetString(tokenProperty)
			if err != nil {
				return err
			}
			client := homewizard.NewAPIv2Client(url, homewizard.WithToken(token))
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
	listUsers.Flags().StringP(tokenProperty, "", "", "The user authentication token")

	cmd.AddCommand(createLocalUser)
	cmd.AddCommand(listUsers)
	return cmd
}

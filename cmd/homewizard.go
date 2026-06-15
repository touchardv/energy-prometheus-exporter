package cmd

import (
	"fmt"
	"os"
	"strings"

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

	read := &cobra.Command{
		Use:   "read",
		Short: "Read current measurements from the Homewizard device",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cmd.Parent().Flags().GetString(flagHomewizardURL)
			if err != nil {
				return err
			}
			if url == "" {
				return fmt.Errorf("homewizard-url must be specified")
			}
			token, err := cmd.Flags().GetString(flagHomewizardToken)
			if err != nil {
				return err
			}
			client := homewizard.NewAuthenticatedAPIv2Client(url, token)
			m, err := client.GetMeasurement()
			if err != nil {
				return err
			}

			rows := []struct {
				name string
				val  string
				unit string
				desc string
			}{
				{"Tariff", fmt.Sprintf("%d", m.Tariff), "", "The active tariff (matches tariff 1 or 2)"},
				{"Timestamp", m.Timestamp, "", "The measurement timestamp (format: YYMMDDhhmmssX)"},
				{"Energy Import (total)", fmt.Sprintf("%.3f", m.EnergyImportkWH), "kWh", "The energy usage meter reading for all tariffs"},
				{"Energy Import T1", fmt.Sprintf("%.3f", m.EnergyImportT1kWH), "kWh", "The energy usage meter reading for tariff 1"},
				{"Energy Import T2", fmt.Sprintf("%.3f", m.EnergyImportT2kWH), "kWh", "The energy usage meter reading for tariff 2"},
				{"Energy Export (total)", fmt.Sprintf("%.3f", m.EnergyExportkWH), "kWh", "The energy feed-in meter reading for all tariffs"},
				{"Energy Export T1", fmt.Sprintf("%.3f", m.EnergyExportT1kWH), "kWh", "The energy feed-in meter reading for tariff 1"},
				{"Energy Export T2", fmt.Sprintf("%.3f", m.EnergyExportT2kWH), "kWh", "The energy feed-in meter reading for tariff 2"},
				{"Power (total)", fmt.Sprintf("%d", m.PowerW), "W", "The total active usage (sum of all phases)"},
				{"Power L1", fmt.Sprintf("%d", m.PowerL1W), "W", "The active usage for phase 1"},
				{"Power L2", fmt.Sprintf("%d", m.PowerL2W), "W", "The active usage for phase 2"},
				{"Power L3", fmt.Sprintf("%d", m.PowerL3W), "W", "The active usage for phase 3"},
				{"Voltage", fmt.Sprintf("%.3f", m.VoltageV), "V", "The active voltage"},
				{"Voltage L1", fmt.Sprintf("%.3f", m.VoltageL1V), "V", "The active voltage for phase 1"},
				{"Voltage L2", fmt.Sprintf("%.3f", m.VoltageL2V), "V", "The active voltage for phase 2"},
				{"Voltage L3", fmt.Sprintf("%.3f", m.VoltageL3V), "V", "The active voltage for phase 3"},
				{"Current", fmt.Sprintf("%.3f", m.CurrentA), "A", "The active current (sum of absolute values)"},
				{"Current L1", fmt.Sprintf("%.3f", m.CurrentL1A), "A", "The active current for phase 1"},
				{"Current L2", fmt.Sprintf("%.3f", m.CurrentL2A), "A", "The active current for phase 2"},
				{"Current L3", fmt.Sprintf("%.3f", m.CurrentL3A), "A", "The active current for phase 3"},
			}

			for _, ext := range m.External {
				rows = append(rows, struct {
					name string
					val  string
					unit string
					desc string
				}{
					name: fmt.Sprintf("External %s", ext.Type),
					val:  fmt.Sprintf("%.3f", ext.Value),
					unit: ext.Unit,
					desc: fmt.Sprintf("External %s measurement", ext.Type),
				})
			}

			maxName := len("Measurement")
			maxVal := len("Value")
			maxUnit := len("Unit")

			for _, r := range rows {
				if len(r.name) > maxName {
					maxName = len(r.name)
				}
				if len(r.val) > maxVal {
					maxVal = len(r.val)
				}
				if len(r.unit) > maxUnit {
					maxUnit = len(r.unit)
				}
			}

			formatStr := fmt.Sprintf("  %%-%ds  %%-%ds  %%-%ds  %%s\n", maxName, maxVal, maxUnit)

			totalWidth := 2 + maxName + 2 + maxVal + 2 + maxUnit + 2 + 50
			if totalWidth < 100 {
				totalWidth = 100
			}
			sep := strings.Repeat("─", totalWidth)

			fmt.Println(sep)
			fmt.Printf(formatStr, "Measurement", "Value", "Unit", "Description")
			fmt.Println(sep)
			for _, r := range rows {
				unitStr := r.unit
				if unitStr == "" {
					unitStr = "-"
				}
				fmt.Printf(formatStr, r.name, r.val, unitStr, r.desc)
			}
			fmt.Println(sep)
			return nil
		},
	}
	read.Flags().StringP(flagHomewizardToken, "", "", "Homewizard user authentication token")

	cmd.AddCommand(createLocalUser)
	cmd.AddCommand(listUsers)
	cmd.AddCommand(read)
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

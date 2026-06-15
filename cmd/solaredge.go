package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/touchardv/energy-prometheus-exporter/pkg/solaredge"
	"github.com/touchardv/energy-prometheus-exporter/pkg/sunspec"
)

const (
	flagSolaredgeHost                    = "solaredge-host"
	flagSolaredgePort                    = "solaredge-port"
	flagSolaredgePollMinimumPollInterval = "solaredge-min-poll-interval"
	flagSolaredgeUnitID                  = "solaredge-unit-id"
	flagSolaredgeMeterSlot               = "solaredge-meter-slot"
	flagSolaredgeBatterySlot             = "solaredge-battery-slot"
)

func addSolaredgeCollectorFlags(fs *pflag.FlagSet) {
	addSolaredgeFlags(fs)
	fs.Duration(flagSolaredgePollMinimumPollInterval, 1*time.Minute, "SolarEdge inverter modbus minimum poll interval duration")
}

func addSolaredgeFlags(fs *pflag.FlagSet) {
	fs.String(flagSolaredgeHost, "", "SolarEdge inverter IP address or hostname")
	fs.Int(flagSolaredgePort, 502, "SolarEdge inverter modbus TCP port")
	fs.Int(flagSolaredgeUnitID, 1, "SolarEdge inverter modbus unit ID")
	fs.Int(flagSolaredgeMeterSlot, 1, "SolarEdge meter slot: 1, 2 or 3")
	fs.Int(flagSolaredgeBatterySlot, 1, "SolarEdge battery slot: 1 or 2")
}

func addSolaredgeCollector(fs *pflag.FlagSet, reg *prometheus.Registry, ctx context.Context) error {
	host, err := fs.GetString(flagSolaredgeHost)
	if err != nil {
		return err
	}
	if host == "" {
		log.Println("solaredge: host not set, collector disabled")
		return nil
	}
	port, err := fs.GetInt(flagSolaredgePort)
	if err != nil {
		return err
	}
	log.Printf("solaredge: collector enabled (host/port: %s:%d)", host, port)
	interval, err := fs.GetDuration(flagSolaredgePollMinimumPollInterval)
	if err != nil {
		return err
	}
	unitID, err := fs.GetInt(flagSolaredgeUnitID)
	if err != nil {
		return err
	}
	meterSlot, err := fs.GetInt(flagSolaredgeMeterSlot)
	if err != nil {
		return err
	}
	batterySlot, err := fs.GetInt(flagSolaredgeBatterySlot)
	if err != nil {
		return err
	}
	collector, err := solaredge.NewCollector(host, port, interval, unitID, meterSlot, batterySlot, ctx)
	if err != nil {
		return err
	}
	reg.MustRegister(collector)
	return nil
}

func newSolaredgeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solaredge",
		Short: "Interact with a SolarEdge inverter",
	}
	addSolaredgeFlags(cmd.PersistentFlags())

	readCmd := &cobra.Command{
		Use:   "read",
		Short: "Read SunSpec models from the SolarEdge inverter and print them",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := cmd.Flags().GetString(flagSolaredgeHost)
			if err != nil {
				return err
			}
			if host == "" {
				return fmt.Errorf("solaredge-host must be specified")
			}
			port, err := cmd.Flags().GetInt(flagSolaredgePort)
			if err != nil {
				return err
			}
			unitID, err := cmd.Flags().GetInt(flagSolaredgeUnitID)
			if err != nil {
				return err
			}
			meterSlot, err := cmd.Flags().GetInt(flagSolaredgeMeterSlot)
			if err != nil {
				return err
			}
			batterySlot, err := cmd.Flags().GetInt(flagSolaredgeBatterySlot)
			if err != nil {
				return err
			}

			url := fmt.Sprintf("tcp://%s:%d", host, port)
			client, err := modbus.NewClient(&modbus.ClientConfiguration{
				URL:     url,
				Timeout: 15 * time.Second,
			})
			if err != nil {
				return fmt.Errorf("creating client for %s: %w", url, err)
			}
			if err := client.Open(); err != nil {
				return fmt.Errorf("connecting to %s: %w", url, err)
			}
			defer client.Close()

			if err := client.SetUnitId(uint8(unitID)); err != nil {
				return fmt.Errorf("setting unit ID: %w", err)
			}

			reader := sunspec.NewReader(client)

			// Read & Print Common Model
			commonModel := sunspec.CommonModel()
			if fields, err := reader.Read(commonModel); err != nil {
				return fmt.Errorf("common model: %w", err)
			} else {
				sunspec.Print(os.Stdout, commonModel, fields)
			}

			// Read & Print Inverter Model
			inverterModel := sunspec.InverterModel()
			if fields, err := reader.Read(inverterModel); err != nil {
				return fmt.Errorf("inverter model: %w", err)
			} else {
				sunspec.Print(os.Stdout, inverterModel, fields)
			}

			// Read & Print Meter Model
			meterModel := sunspec.MeterModel(meterSlot)
			if fields, err := reader.Read(meterModel); err != nil {
				return fmt.Errorf("meter model: %w", err)
			} else {
				sunspec.Print(os.Stdout, meterModel, fields)
			}

			// Read & Print Battery Info Model
			battInfoModel := sunspec.BatteryInfoModel(batterySlot)
			if fields, err := reader.Read(battInfoModel); err != nil {
				log.Printf("solaredge: could not read battery info model: %v", err)
			} else {
				sunspec.Print(os.Stdout, battInfoModel, fields)

				// Read & Print Battery Data Model (only try if battery info succeeded)
				battDataModel := sunspec.BatteryDataModel(batterySlot)
				if fields, err := reader.Read(battDataModel); err != nil {
					log.Printf("solaredge: could not read battery data model: %v", err)
				} else {
					sunspec.Print(os.Stdout, battDataModel, fields)
				}
			}

			return nil
		},
	}

	cmd.AddCommand(readCmd)
	return cmd
}

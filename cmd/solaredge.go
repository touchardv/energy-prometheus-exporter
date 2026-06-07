package cmd

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/touchardv/energy-prometheus-exporter/pkg/solaredge"
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
	fs.String(flagSolaredgeHost, "", "SolarEdge inverter IP address or hostname")
	fs.Int(flagSolaredgePort, 502, "SolarEdge inverter modbus TCP port")
	fs.Duration(flagSolaredgePollMinimumPollInterval, 1*time.Minute, "SolarEdge inverter modbus minimum poll interval duration")
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

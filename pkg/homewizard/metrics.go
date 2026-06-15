package homewizard

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// See https://api-documentation.homewizard.com/docs/v2/measurement
var (
	tariffDesc = prometheus.NewDesc(
		"homewizard_p1_tariff",
		"The active tariff, matches one of the totals.",
		[]string{},
		nil,
	)
	energyImportKWHDesc = prometheus.NewDesc(
		"homewizard_energy_import_kwh",
		"The energy usage meter reading for all tariffs in kWh.",
		[]string{},
		nil,
	)
	energyImportT1KWHDesc = prometheus.NewDesc(
		"homewizard_energy_import_t1_kwh",
		"The energy usage meter reading for tariff 1 in kWh.",
		[]string{},
		nil,
	)
	energyImportT2KWHDesc = prometheus.NewDesc(
		"homewizard_energy_import_t2_kwh",
		"The energy usage meter reading for tariff 2 in kWh.",
		[]string{},
		nil,
	)
	energyExportKWHDesc = prometheus.NewDesc(
		"homewizard_energy_export_kwh",
		"The energy feed-in meter reading for all tariffs in kWh.",
		[]string{},
		nil,
	)
	energyExportT1KWHDesc = prometheus.NewDesc(
		"homewizard_energy_export_t1_kwh",
		"The energy feed-in meter reading for tariff 1 in kWh.",
		[]string{},
		nil,
	)
	energyExportT2KWHDesc = prometheus.NewDesc(
		"homewizard_energy_export_t2_kwh",
		"The energy feed-in meter reading for tariff 2 in kWh.",
		[]string{},
		nil,
	)
	powerWDesc = prometheus.NewDesc(
		"homewizard_power_w",
		"The total active usage in watt, this value is the sum of all phases, so if l1=200, l2=300, l3=-100 then this value is 400.",
		[]string{},
		nil,
	)
	powerL1WDesc = prometheus.NewDesc(
		"homewizard_power_l1_w",
		"The active usage for phase 1 in watt, will be negative when exporting.",
		[]string{},
		nil,
	)
	powerL2WDesc = prometheus.NewDesc(
		"homewizard_power_l2_w",
		"The active usage for phase 2 in watt, will be negative when exporting.",
		[]string{},
		nil,
	)
	powerL3WDesc = prometheus.NewDesc(
		"homewizard_power_l3_w",
		"The active usage for phase 3 in watt, will be negative when exporting.",
		[]string{},
		nil,
	)
	voltageVDesc = prometheus.NewDesc(
		"homewizard_voltage_v",
		"The active voltage in volt.",
		[]string{},
		nil,
	)
	voltageL1VDesc = prometheus.NewDesc(
		"homewizard_voltage_l1_v",
		"The active voltage for phase 1 in volt.",
		[]string{},
		nil,
	)
	voltageL2VDesc = prometheus.NewDesc(
		"homewizard_voltage_l2_v",
		"The active voltage for phase 2 in volt.",
		[]string{},
		nil,
	)
	voltageL3VDesc = prometheus.NewDesc(
		"homewizard_voltage_l3_v",
		"The active voltage for phase 3 in volt.",
		[]string{},
		nil,
	)
	currentADesc = prometheus.NewDesc(
		"homewizard_current_a",
		"The active current in ampere, this value is the sum of absolute values, so if l1=2, l2=3, l3=-1 then this value is 6.",
		[]string{},
		nil,
	)
	currentL1ADesc = prometheus.NewDesc(
		"homewizard_current_l1_a",
		"The active current for phase 1 in ampere, will be negative when exporting.",
		[]string{},
		nil,
	)
	currentL2ADesc = prometheus.NewDesc(
		"homewizard_current_l2_a",
		"The active current for phase 2 in ampere, will be negative when exporting.",
		[]string{},
		nil,
	)
	currentL3ADesc = prometheus.NewDesc(
		"homewizard_current_l3_a",
		"The active current for phase 3 in ampere, will be negative when exporting.",
		[]string{},
		nil,
	)
	extGasMeterDesc = prometheus.NewDesc(
		"homewizard_ext_gas_meter_m3",
		"The total of gaz in m3.",
		[]string{},
		nil,
	)
)

// Collect implements [prometheus.Collector].
func (c APIv2Client) Collect(ch chan<- prometheus.Metric) {
	logrus.Debug("Collecting metrics...")
	m, err := c.GetMeasurement()
	if err != nil {
		logrus.Error(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(tariffDesc, prometheus.GaugeValue, float64(m.Tariff))
	ch <- prometheus.MustNewConstMetric(energyImportKWHDesc, prometheus.CounterValue, m.EnergyImportkWH)
	ch <- prometheus.MustNewConstMetric(energyImportT1KWHDesc, prometheus.CounterValue, m.EnergyImportT1kWH)
	ch <- prometheus.MustNewConstMetric(energyImportT2KWHDesc, prometheus.CounterValue, m.EnergyImportT2kWH)
	ch <- prometheus.MustNewConstMetric(energyExportKWHDesc, prometheus.CounterValue, m.EnergyExportkWH)
	ch <- prometheus.MustNewConstMetric(energyExportT1KWHDesc, prometheus.CounterValue, m.EnergyExportT1kWH)
	ch <- prometheus.MustNewConstMetric(energyExportT2KWHDesc, prometheus.CounterValue, m.EnergyExportT2kWH)
	ch <- prometheus.MustNewConstMetric(powerWDesc, prometheus.GaugeValue, float64(m.PowerW))
	ch <- prometheus.MustNewConstMetric(powerL1WDesc, prometheus.GaugeValue, float64(m.PowerL1W))
	ch <- prometheus.MustNewConstMetric(powerL2WDesc, prometheus.GaugeValue, float64(m.PowerL2W))
	ch <- prometheus.MustNewConstMetric(powerL3WDesc, prometheus.GaugeValue, float64(m.PowerL3W))
	ch <- prometheus.MustNewConstMetric(voltageVDesc, prometheus.GaugeValue, m.VoltageV)
	ch <- prometheus.MustNewConstMetric(voltageL1VDesc, prometheus.GaugeValue, m.VoltageL1V)
	ch <- prometheus.MustNewConstMetric(voltageL2VDesc, prometheus.GaugeValue, m.VoltageL2V)
	ch <- prometheus.MustNewConstMetric(voltageL3VDesc, prometheus.GaugeValue, m.VoltageL3V)
	ch <- prometheus.MustNewConstMetric(currentADesc, prometheus.GaugeValue, m.CurrentA)
	ch <- prometheus.MustNewConstMetric(currentL1ADesc, prometheus.GaugeValue, m.CurrentL1A)
	ch <- prometheus.MustNewConstMetric(currentL2ADesc, prometheus.GaugeValue, m.CurrentL2A)
	ch <- prometheus.MustNewConstMetric(currentL3ADesc, prometheus.GaugeValue, m.CurrentL3A)
	for _, em := range m.External {
		if em.Type == "gas_meter" {
			ch <- prometheus.MustNewConstMetric(extGasMeterDesc, prometheus.GaugeValue, em.Value)
		}
	}
}

// Describe implements [prometheus.Collector].
func (c APIv2Client) Describe(ch chan<- *prometheus.Desc) {
	logrus.Debug("Collecting metric descriptions...")
	ch <- tariffDesc
	ch <- energyImportKWHDesc
	ch <- energyImportT1KWHDesc
	ch <- energyImportT2KWHDesc
	ch <- energyExportKWHDesc
	ch <- energyExportT1KWHDesc
	ch <- energyExportT2KWHDesc
	ch <- powerWDesc
	ch <- powerL1WDesc
	ch <- powerL2WDesc
	ch <- powerL3WDesc
	ch <- voltageVDesc
	ch <- voltageL1VDesc
	ch <- voltageL2VDesc
	ch <- voltageL3VDesc
	ch <- currentADesc
	ch <- currentL1ADesc
	ch <- currentL2ADesc
	ch <- currentL3ADesc
	ch <- extGasMeterDesc
}

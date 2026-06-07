package solaredge

import "github.com/prometheus/client_golang/prometheus"

// See
// https://knowledge-center.solaredge.com/sites/kc/files/sunspec-implementation-technical-note.pdf
// https://sunspec.org/wp-content/uploads/2009/03/SunSpec-Alliance-Specification-Energy-Storage-ModelsD4rev0.pdf
var (
	descUp = prometheus.NewDesc(
		"solaredge_up", "1 if the last poll succeeded", nil, nil)

	descInverterACPower = prometheus.NewDesc(
		"solaredge_inverter_ac_power_watts", "AC output power", nil, nil)
	descInverterDCPower = prometheus.NewDesc(
		"solaredge_inverter_dc_power_watts", "DC input power from panels", nil, nil)
	descInverterACEnergy = prometheus.NewDesc(
		"solaredge_inverter_ac_energy_wh_total", "AC lifetime energy production", nil, nil)
	descInverterACCurrent = prometheus.NewDesc(
		"solaredge_inverter_ac_current_amps", "AC total current", nil, nil)
	descInverterACVoltage = prometheus.NewDesc(
		"solaredge_inverter_ac_voltage_an_volts", "AC voltage phase A-N", nil, nil)
	descInverterACFrequency = prometheus.NewDesc(
		"solaredge_inverter_ac_frequency_hz", "AC frequency", nil, nil)
	descInverterACPowerFactor = prometheus.NewDesc(
		"solaredge_inverter_ac_power_factor_pct", "AC power factor", nil, nil)
	descInverterDCVoltage = prometheus.NewDesc(
		"solaredge_inverter_dc_voltage_volts", "DC voltage", nil, nil)
	descInverterDCCurrent = prometheus.NewDesc(
		"solaredge_inverter_dc_current_amps", "DC current", nil, nil)
	descInverterTempSink = prometheus.NewDesc(
		"solaredge_inverter_temp_heatsink_celsius", "Heat sink temperature", nil, nil)
	descInverterStatus = prometheus.NewDesc(
		"solaredge_inverter_status",
		"Operating state (1=off 2=sleeping 3=starting 4=mppt 5=throttled 6=shutdown 7=fault 8=standby)",
		nil, nil)

	descMeterACPower = prometheus.NewDesc(
		"solaredge_meter_ac_power_watts",
		"Total real power (positive=export, negative=import)", nil, nil)
	descMeterACPowerPhase = prometheus.NewDesc(
		"solaredge_meter_ac_power_phase_watts", "Real power per phase",
		[]string{"phase"}, nil)
	descMeterACCurrent = prometheus.NewDesc(
		"solaredge_meter_ac_current_amps", "AC total current", nil, nil)
	descMeterACVoltage = prometheus.NewDesc(
		"solaredge_meter_ac_voltage_an_volts", "AC voltage line-neutral average", nil, nil)
	descMeterACFrequency = prometheus.NewDesc(
		"solaredge_meter_ac_frequency_hz", "AC frequency", nil, nil)
	descMeterExported = prometheus.NewDesc(
		"solaredge_meter_exported_wh_total", "Total exported real energy", nil, nil)
	descMeterImported = prometheus.NewDesc(
		"solaredge_meter_imported_wh_total", "Total imported real energy", nil, nil)
	descMeterACPowerFactor = prometheus.NewDesc(
		"solaredge_meter_ac_power_factor_pct", "Average power factor", nil, nil)

	descBatterySoC = prometheus.NewDesc(
		"solaredge_battery_soc_pct", "Battery state of charge in percent", nil, nil)
	descBatterySoH = prometheus.NewDesc(
		"solaredge_battery_soh_pct", "Battery state of health in percent", nil, nil)
	descBatteryPower = prometheus.NewDesc(
		"solaredge_battery_power_watts",
		"Battery power in watts (positive=charging, negative=discharging)", nil, nil)
	descBatteryVoltage = prometheus.NewDesc(
		"solaredge_battery_voltage_volts", "Battery voltage in volts", nil, nil)
	descBatteryCurrent = prometheus.NewDesc(
		"solaredge_battery_current_amps",
		"Battery current in amps (positive=charging, negative=discharging)", nil, nil)
	descBatteryTempAvg = prometheus.NewDesc(
		"solaredge_battery_temp_avg_celsius", "Battery average temperature in Celsius", nil, nil)
	descBatteryTempMax = prometheus.NewDesc(
		"solaredge_battery_temp_max_celsius", "Battery maximum temperature in Celsius", nil, nil)
	descBatteryEnergyAvailable = prometheus.NewDesc(
		"solaredge_battery_energy_available_wh", "Battery available energy in watt-hours", nil, nil)
	descBatteryStatus = prometheus.NewDesc(
		"solaredge_battery_status",
		"Battery operating status (0=Off 1=Standby 2=Init 3=Charging 4=Discharging 5=Fault 6=Idle)", nil, nil)
	descBatteryExported = prometheus.NewDesc(
		"solaredge_battery_exported_wh_total",
		"Total energy discharged from battery in watt-hours", nil, nil)
	descBatteryImported = prometheus.NewDesc(
		"solaredge_battery_imported_wh_total",
		"Total energy charged into battery in watt-hours", nil, nil)
)

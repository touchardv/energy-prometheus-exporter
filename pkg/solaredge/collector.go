package solaredge

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
	"github.com/touchardv/energy-prometheus-exporter/pkg/sunspec"
)

type snapshot struct {
	inverter map[string]sunspec.DecodedField
	meter    map[string]sunspec.DecodedField
	battery  map[string]sunspec.DecodedField // nil when battery is absent
	ok       bool                            // false if the last poll failed
}

// Collector implements prometheus.Collector. It maintains a persistent Modbus
// TCP connection and fetches metrics on-demand, caching the results for a
// minimum of interval to avoid overloading the inverter.
type Collector struct {
	host        string
	port        int
	minimumPoll time.Duration
	unitID      int
	meterSlot   int
	batterySlot int

	mu           sync.Mutex
	cur          snapshot
	lastPollTime time.Time
	client       *modbus.ModbusClient
	reader       *sunspec.Reader
}

// NewCollector constructs a Collector.
// A background goroutine monitors ctx to close any active Modbus client
// connection when the context is cancelled.
func NewCollector(host string, port int, minimumPoll time.Duration, unitID int, meterSlot int, batterySlot int, ctx context.Context) (*Collector, error) {

	c := &Collector{
		host:        host,
		port:        port,
		minimumPoll: minimumPoll,
		unitID:      unitID,
		meterSlot:   meterSlot,
		batterySlot: batterySlot,
	}

	go func() {
		<-ctx.Done()
		c.mu.Lock()
		c.closeClient()
		c.mu.Unlock()
	}()

	return c, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range allDescs() {
		ch <- d
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.lastPollTime.IsZero() && time.Since(c.lastPollTime) < c.minimumPoll {
		c.emitMetrics(ch)
		return
	}

	if err := c.ensureConnected(); err != nil {
		log.Printf("solaredge: connect error: %v", err)
		c.cur = snapshot{ok: false}
		c.lastPollTime = time.Now()
		c.emitMetrics(ch)
		return
	}

	if err := c.poll(); err != nil {
		log.Printf("solaredge: poll error: %v", err)
		c.cur = snapshot{ok: false}
		c.closeClient()
		c.lastPollTime = time.Now()
		c.emitMetrics(ch)
		return
	}

	c.lastPollTime = time.Now()
	c.emitMetrics(ch)
}

// ensureConnected makes sure c.client is initialized and connected.
// If it is not, a connection is dialed.
// Caller must hold c.mu.
func (c *Collector) ensureConnected() error {
	if c.client != nil {
		return nil
	}

	client, reader, err := c.dial()
	if err != nil {
		return err
	}

	log.Printf("solaredge: connected to %s:%d", c.host, c.port)

	logDeviceInfo(reader)
	logMeterInfo(reader, c.meterSlot)
	logBatteryInfo(reader, c.batterySlot)

	c.client = client
	c.reader = reader
	return nil
}

func (c *Collector) dial() (*modbus.ModbusClient, *sunspec.Reader, error) {
	url := fmt.Sprintf("tcp://%s:%d", c.host, c.port)

	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     url,
		Timeout: 15 * time.Second,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("creating client for %s: %w", url, err)
	}

	if err := client.Open(); err != nil {
		return nil, nil, fmt.Errorf("connecting to %s: %w", url, err)
	}

	if err := client.SetUnitId(uint8(c.unitID)); err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("setting unit ID: %w", err)
	}

	return client, sunspec.NewReader(client), nil
}

// closeClient closes the active client connection if open.
// Caller must hold c.mu.
func (c *Collector) closeClient() {
	if c.client != nil {
		c.client.Close()
		c.client = nil
		c.reader = nil
	}
}

// isTransactionIDError returns true when err is a Modbus transaction ID
// mismatch. This happens when the inverter has stale buffered responses from
// a previous session. The correct response is to reconnect immediately rather
// than waiting a full interval — the next connection will be clean.
func isTransactionIDError(err error) bool {
	if err == nil {
		return false
	}
	return err == modbus.ErrBadTransactionId ||
		strings.Contains(err.Error(), "transaction id")
}

// poll reads both models and updates the snapshot atomically.
// Caller must hold c.mu.
func (c *Collector) poll() error {
	invFields, err := c.reader.Read(sunspec.InverterModel())
	if err != nil {
		if isTransactionIDError(err) {
			log.Printf("solaredge: stale transaction ID on inverter read — reconnecting")
		}
		return fmt.Errorf("inverter: %w", err)
	}

	meterFields, err := c.reader.Read(sunspec.MeterModel(c.meterSlot))
	if err != nil {
		if isTransactionIDError(err) {
			log.Printf("solaredge: stale transaction ID on meter read — reconnecting")
		}
		return fmt.Errorf("meter: %w", err)
	}

	snap := snapshot{
		inverter: sunspec.FieldMap(invFields),
		meter:    sunspec.FieldMap(meterFields),
		ok:       true,
	}

	// Battery is non-fatal — not all installations have one.
	battFields, err := c.reader.Read(sunspec.BatteryDataModel(c.batterySlot))
	if err != nil {
		if !isTransactionIDError(err) {
			log.Printf("solaredge: battery poll error: %v", err)
		}
	} else {
		snap.battery = sunspec.FieldMap(battFields)
	}

	c.cur = snap
	return nil
}

func (c *Collector) emitMetrics(ch chan<- prometheus.Metric) {
	s := c.cur

	up := 0.0
	if s.ok {
		up = 1.0
	}
	ch <- prometheus.MustNewConstMetric(descUp, prometheus.GaugeValue, up)

	if !s.ok {
		return
	}

	emitInverter(s.inverter, ch)
	emitMeter(s.meter, ch)
	if s.battery != nil {
		emitBattery(s.battery, ch)
	}
}

func emitInverter(fm map[string]sunspec.DecodedField, ch chan<- prometheus.Metric) {
	gauge := func(desc *prometheus.Desc, name string) {
		emit(ch, desc, prometheus.GaugeValue, fm, name)
	}
	gauge(descInverterACPower, "I_AC_Power")
	gauge(descInverterDCPower, "I_DC_Power")
	gauge(descInverterACCurrent, "I_AC_Current")
	gauge(descInverterACVoltage, "I_AC_VoltageAN")
	gauge(descInverterACFrequency, "I_AC_Frequency")
	gauge(descInverterACPowerFactor, "I_AC_PF")
	gauge(descInverterDCVoltage, "I_DC_Voltage")
	gauge(descInverterDCCurrent, "I_DC_Current")
	gauge(descInverterTempSink, "I_Temp_Sink")
	gauge(descInverterStatus, "I_Status")
	emit(ch, descInverterACEnergy, prometheus.CounterValue, fm, "I_AC_Energy_WH")
}

func emitMeter(fm map[string]sunspec.DecodedField, ch chan<- prometheus.Metric) {
	gauge := func(desc *prometheus.Desc, name string, labels ...string) {
		emit(ch, desc, prometheus.GaugeValue, fm, name, labels...)
	}
	gauge(descMeterACPower, "M_AC_Power")
	gauge(descMeterACCurrent, "M_AC_Current")
	gauge(descMeterACVoltage, "M_AC_Voltage_LN")
	gauge(descMeterACFrequency, "M_AC_Freq")
	gauge(descMeterACPowerFactor, "M_AC_PF")
	gauge(descMeterACPowerPhase, "M_AC_Power_A", "A")
	gauge(descMeterACPowerPhase, "M_AC_Power_B", "B")
	gauge(descMeterACPowerPhase, "M_AC_Power_C", "C")
	emit(ch, descMeterExported, prometheus.CounterValue, fm, "M_Exported")
	emit(ch, descMeterImported, prometheus.CounterValue, fm, "M_Imported")
}

func emitBattery(fm map[string]sunspec.DecodedField, ch chan<- prometheus.Metric) {
	gauge := func(desc *prometheus.Desc, name string) {
		emit(ch, desc, prometheus.GaugeValue, fm, name)
	}
	gauge(descBatterySoC, "B_SoC")
	gauge(descBatterySoH, "B_SoH")
	gauge(descBatteryPower, "B_Power")
	gauge(descBatteryVoltage, "B_Voltage")
	gauge(descBatteryCurrent, "B_Current")
	gauge(descBatteryTempAvg, "B_TempAverage")
	gauge(descBatteryTempMax, "B_TempMax")
	gauge(descBatteryEnergyAvailable, "B_EnergyAvailable")
	gauge(descBatteryStatus, "B_Status")
	emit(ch, descBatteryExported, prometheus.CounterValue, fm, "B_ExportedEnergy")
	emit(ch, descBatteryImported, prometheus.CounterValue, fm, "B_ImportedEnergy")
}

func emit(ch chan<- prometheus.Metric, desc *prometheus.Desc, vt prometheus.ValueType,
	fm map[string]sunspec.DecodedField, name string, labels ...string) {
	f, ok := fm[name]
	if !ok || f.NotImplemented {
		return
	}
	ch <- prometheus.MustNewConstMetric(desc, vt, f.Value, labels...)
}

func allDescs() []*prometheus.Desc {
	return []*prometheus.Desc{
		descUp,
		descInverterACPower, descInverterDCPower, descInverterACEnergy,
		descInverterACCurrent, descInverterACVoltage, descInverterACFrequency,
		descInverterACPowerFactor, descInverterDCVoltage, descInverterDCCurrent,
		descInverterTempSink, descInverterStatus,
		descMeterACPower, descMeterACPowerPhase, descMeterACCurrent,
		descMeterACVoltage, descMeterACFrequency, descMeterExported,
		descMeterImported, descMeterACPowerFactor,
		descBatterySoC, descBatterySoH, descBatteryPower,
		descBatteryVoltage, descBatteryCurrent, descBatteryTempAvg,
		descBatteryTempMax, descBatteryEnergyAvailable, descBatteryStatus,
		descBatteryExported, descBatteryImported,
	}
}

// logDeviceInfo reads the SunSpec Common block and logs the inverter identity.
// Called once after each successful Modbus connect. Failures are non-fatal.
func logDeviceInfo(reader *sunspec.Reader) {
	fields, err := reader.Read(sunspec.CommonModel())
	if err != nil {
		log.Printf("solaredge: could not read device info: %v", err)
		return
	}
	fm := sunspec.FieldMap(fields)
	str := func(name string) string {
		if f, ok := fm[name]; ok && f.IsStr {
			return f.StrVal
		}
		return "N/A"
	}
	log.Printf("solaredge: inverter manufacturer  : %s", str("C_Manufacturer"))
	log.Printf("solaredge: inverter model         : %s", str("C_Model"))
	log.Printf("solaredge: inverter serial number : %s", str("C_SerialNumber"))
	log.Printf("solaredge: inverter firmware      : %s", str("C_Version"))
}

// logMeterInfo reads the meter Common block (embedded at the start of the
// meter model block) and logs its identity. Failures are non-fatal.
func logMeterInfo(reader *sunspec.Reader, slot int) {
	fields, err := reader.Read(sunspec.MeterModel(slot))
	if err != nil {
		log.Printf("solaredge: could not read meter info: %v", err)
		return
	}
	fm := sunspec.FieldMap(fields)
	str := func(name string) string {
		if f, ok := fm[name]; ok && f.IsStr {
			return f.StrVal
		}
		return "N/A"
	}
	log.Printf("solaredge: meter manufacturer  : %s", str("C_Manufacturer"))
	log.Printf("solaredge: meter model         : %s", str("C_Model"))
	log.Printf("solaredge: meter serial number : %s", str("C_SerialNumber"))
	log.Printf("solaredge: meter option        : %s", str("C_Option"))
	log.Printf("solaredge: meter firmware      : %s", str("C_Version"))
}

// logBatteryInfo reads the battery nameplate block and logs its identity.
// Failures are non-fatal — not all installations have a battery.
func logBatteryInfo(reader *sunspec.Reader, slot int) {
	fields, err := reader.Read(sunspec.BatteryInfoModel(slot))
	if err != nil {
		log.Printf("solaredge: battery not available (slot %d): %v", slot, err)
		return
	}
	fm := sunspec.FieldMap(fields)

	str := func(name string) string {
		if f, ok := fm[name]; ok && f.IsStr && !f.NotImplemented {
			return f.StrVal
		}
		return "N/A"
	}
	floatVal := func(name string) string {
		if f, ok := fm[name]; ok && !f.NotImplemented {
			return fmt.Sprintf("%.0f", f.Value)
		}
		return "N/A"
	}

	log.Printf("solaredge: battery manufacturer  : %s", str("B_Manufacturer"))
	log.Printf("solaredge: battery model         : %s", str("B_Model"))
	log.Printf("solaredge: battery serial number : %s", str("B_SerialNumber"))
	log.Printf("solaredge: battery firmware      : %s", str("B_Version"))
	log.Printf("solaredge: battery rated energy  : %s Wh", floatVal("B_RatedEnergy"))
	log.Printf("solaredge: battery max charge    : %s W", floatVal("B_MaxChargePower"))
	log.Printf("solaredge: battery max discharge : %s W", floatVal("B_MaxDischargePower"))
}

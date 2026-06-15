package sunspec

// ─────────────────────────────────────────────────────────────────
// Pre-built model definitions derived from the SolarEdge SunSpec
// Technical Note (v2.5, November 2022).
//
// All offsets are 0-based within the model block.
// Absolute Modbus address = model.BaseAddress + field.Offset
// ─────────────────────────────────────────────────────────────────

const (
	CommonBaseAddress   uint16 = 40000
	InverterBaseAddress uint16 = 40069
)

// CommonModel returns the SunSpec Common Block definition.
// Default base address: 40000 (base-0).
func CommonModel() *Model {
	return &Model{
		Name:        "Common",
		BaseAddress: CommonBaseAddress,
		Fields: []RegisterDef{
			{0, 2, "C_SunSpec_ID", TypeUint32, "", "SunSpec identifier (0x53756e53 = 'SunS')", "", 0},
			{2, 1, "C_SunSpec_DID", TypeUint16, "", "Common Model Block ID (1)", "", 0},
			{3, 1, "C_SunSpec_Length", TypeUint16, "registers", "Block length", "", 0},
			{4, 16, "C_Manufacturer", TypeString, "", "Manufacturer", "", 32},
			{20, 16, "C_Model", TypeString, "", "Device model", "", 32},
			{44, 8, "C_Version", TypeString, "", "Firmware version", "", 16},
			{52, 16, "C_SerialNumber", TypeString, "", "Serial number", "", 32},
			{68, 1, "C_DeviceAddress", TypeUint16, "", "Modbus Unit ID", "", 0},
		},
	}
}

// InverterModel returns the SunSpec Inverter Model block (IDs 101/102/103).
// Default base address: 40069 (base-0).
func InverterModel() *Model {
	return &Model{
		Name:        "Inverter",
		BaseAddress: InverterBaseAddress,
		Fields: []RegisterDef{
			{0, 1, "C_SunSpec_DID", TypeUint16, "", "101=single, 102=split, 103=three-phase", "", 0},
			{1, 1, "C_SunSpec_Length", TypeUint16, "registers", "Model block length (50)", "", 0},
			{2, 1, "I_AC_Current", TypeUint16, "A", "AC Total Current", "I_AC_Current_SF", 0},
			{3, 1, "I_AC_CurrentA", TypeUint16, "A", "Phase A Current", "I_AC_Current_SF", 0},
			{4, 1, "I_AC_CurrentB", TypeUint16, "A", "Phase B Current", "I_AC_Current_SF", 0},
			{5, 1, "I_AC_CurrentC", TypeUint16, "A", "Phase C Current", "I_AC_Current_SF", 0},
			{6, 1, "I_AC_Current_SF", TypeSF, "", "AC Current scale factor", "", 0},
			{7, 1, "I_AC_VoltageAB", TypeUint16, "V", "AC Voltage Phase AB", "I_AC_Voltage_SF", 0},
			{8, 1, "I_AC_VoltageBC", TypeUint16, "V", "AC Voltage Phase BC", "I_AC_Voltage_SF", 0},
			{9, 1, "I_AC_VoltageCA", TypeUint16, "V", "AC Voltage Phase CA", "I_AC_Voltage_SF", 0},
			{10, 1, "I_AC_VoltageAN", TypeUint16, "V", "AC Voltage Phase A-N", "I_AC_Voltage_SF", 0},
			{11, 1, "I_AC_VoltageBN", TypeUint16, "V", "AC Voltage Phase B-N", "I_AC_Voltage_SF", 0},
			{12, 1, "I_AC_VoltageCN", TypeUint16, "V", "AC Voltage Phase C-N", "I_AC_Voltage_SF", 0},
			{13, 1, "I_AC_Voltage_SF", TypeSF, "", "AC Voltage scale factor", "", 0},
			{14, 1, "I_AC_Power", TypeInt16, "W", "AC Power", "I_AC_Power_SF", 0},
			{15, 1, "I_AC_Power_SF", TypeSF, "", "AC Power scale factor", "", 0},
			{16, 1, "I_AC_Frequency", TypeUint16, "Hz", "AC Frequency", "I_AC_Frequency_SF", 0},
			{17, 1, "I_AC_Frequency_SF", TypeSF, "", "AC Frequency scale factor", "", 0},
			{18, 1, "I_AC_VA", TypeInt16, "VA", "Apparent Power", "I_AC_VA_SF", 0},
			{19, 1, "I_AC_VA_SF", TypeSF, "", "Apparent Power scale factor", "", 0},
			{20, 1, "I_AC_VAR", TypeInt16, "VAR", "Reactive Power", "I_AC_VAR_SF", 0},
			{21, 1, "I_AC_VAR_SF", TypeSF, "", "Reactive Power scale factor", "", 0},
			{22, 1, "I_AC_PF", TypeInt16, "%", "Power Factor", "I_AC_PF_SF", 0},
			{23, 1, "I_AC_PF_SF", TypeSF, "", "Power Factor scale factor", "", 0},
			{24, 2, "I_AC_Energy_WH", TypeAcc32, "Wh", "AC Lifetime Energy", "I_AC_Energy_WH_SF", 0},
			{26, 1, "I_AC_Energy_WH_SF", TypeSF, "", "AC Energy scale factor", "", 0},
			{27, 1, "I_DC_Current", TypeUint16, "A", "DC Current", "I_DC_Current_SF", 0},
			{28, 1, "I_DC_Current_SF", TypeSF, "", "DC Current scale factor", "", 0},
			{29, 1, "I_DC_Voltage", TypeUint16, "V", "DC Voltage", "I_DC_Voltage_SF", 0},
			{30, 1, "I_DC_Voltage_SF", TypeSF, "", "DC Voltage scale factor", "", 0},
			{31, 1, "I_DC_Power", TypeInt16, "W", "DC Power", "I_DC_Power_SF", 0},
			{32, 1, "I_DC_Power_SF", TypeSF, "", "DC Power scale factor", "", 0},
			{34, 1, "I_Temp_Sink", TypeInt16, "°C", "Heat Sink Temperature", "I_Temp_SF", 0},
			{37, 1, "I_Temp_SF", TypeSF, "", "Temperature scale factor", "", 0},
			{38, 1, "I_Status", TypeUint16, "", "Operating State", "", 0},
			{39, 1, "I_Status_Vendor", TypeUint16, "", "Vendor Operating State", "", 0},
		},
	}
}

// MeterModel returns the SunSpec Meter Model block (IDs 201–204).
// The base address depends on which meter slot and inverter type:
//
//	Meter 1 (standard):  40121 (base-0)
//	Meter 2 (standard):  40295 (base-0)
//	Meter 3 (standard):  40469 (base-0)
//
// Pass the appropriate base address for your installation.
func MeterModel(slot int) *Model {
	return &Model{
		Name:        "Meter",
		BaseAddress: meterBase(slot),
		Fields: []RegisterDef{
			// ── Common Block ─────────────────────────────────────────────────
			{0, 1, "C_SunSpec_DID", TypeUint16, "", "Common Model Block ID (1)", "", 0},
			{1, 1, "C_SunSpec_Length", TypeUint16, "registers", "Common block length (65)", "", 0},
			{2, 16, "C_Manufacturer", TypeString, "", "Meter manufacturer", "", 32},
			{18, 16, "C_Model", TypeString, "", "Meter model", "", 32},
			{34, 8, "C_Option", TypeString, "", "Export+Import / Production / Consumption", "", 16},
			{42, 8, "C_Version", TypeString, "", "Meter firmware version", "", 16},
			{50, 16, "C_SerialNumber", TypeString, "", "Meter serial number", "", 32},
			{66, 1, "C_DeviceAddress", TypeUint16, "", "Inverter Modbus ID", "", 0},
			// ── Identification ────────────────────────────────────────────────
			{67, 1, "M_SunSpec_DID", TypeUint16, "", "201=1ph 202=split 203=WYE 204=delta", "", 0},
			{68, 1, "M_SunSpec_Length", TypeUint16, "registers", "Meter model block length", "", 0},
			// ── Current ───────────────────────────────────────────────────────
			{69, 1, "M_AC_Current", TypeInt16, "A", "AC Current (sum of active phases)", "M_AC_Current_SF", 0},
			{70, 1, "M_AC_Current_A", TypeInt16, "A", "Phase A AC Current", "M_AC_Current_SF", 0},
			{71, 1, "M_AC_Current_B", TypeInt16, "A", "Phase B AC Current", "M_AC_Current_SF", 0},
			{72, 1, "M_AC_Current_C", TypeInt16, "A", "Phase C AC Current", "M_AC_Current_SF", 0},
			{73, 1, "M_AC_Current_SF", TypeSF, "", "AC Current scale factor", "", 0},
			// ── Voltage (L-N) ─────────────────────────────────────────────────
			{74, 1, "M_AC_Voltage_LN", TypeInt16, "V", "Line-Neutral Voltage (avg)", "M_AC_Voltage_SF", 0},
			{75, 1, "M_AC_Voltage_AN", TypeInt16, "V", "Phase A-N Voltage", "M_AC_Voltage_SF", 0},
			{76, 1, "M_AC_Voltage_BN", TypeInt16, "V", "Phase B-N Voltage", "M_AC_Voltage_SF", 0},
			{77, 1, "M_AC_Voltage_CN", TypeInt16, "V", "Phase C-N Voltage", "M_AC_Voltage_SF", 0},
			// ── Voltage (L-L) ─────────────────────────────────────────────────
			{78, 1, "M_AC_Voltage_LL", TypeInt16, "V", "Line-Line Voltage (avg)", "M_AC_Voltage_SF", 0},
			{79, 1, "M_AC_Voltage_AB", TypeInt16, "V", "Phase A-B Voltage", "M_AC_Voltage_SF", 0},
			{80, 1, "M_AC_Voltage_BC", TypeInt16, "V", "Phase B-C Voltage", "M_AC_Voltage_SF", 0},
			{81, 1, "M_AC_Voltage_CA", TypeInt16, "V", "Phase C-A Voltage", "M_AC_Voltage_SF", 0},
			{82, 1, "M_AC_Voltage_SF", TypeSF, "", "AC Voltage scale factor", "", 0},
			// ── Frequency ─────────────────────────────────────────────────────
			{83, 1, "M_AC_Freq", TypeInt16, "Hz", "AC Frequency", "M_AC_Freq_SF", 0},
			{84, 1, "M_AC_Freq_SF", TypeSF, "", "AC Frequency scale factor", "", 0},
			// ── Real Power ────────────────────────────────────────────────────
			{85, 1, "M_AC_Power", TypeInt16, "W", "Total Real Power", "M_AC_Power_SF", 0},
			{86, 1, "M_AC_Power_A", TypeInt16, "W", "Phase A Real Power", "M_AC_Power_SF", 0},
			{87, 1, "M_AC_Power_B", TypeInt16, "W", "Phase B Real Power", "M_AC_Power_SF", 0},
			{88, 1, "M_AC_Power_C", TypeInt16, "W", "Phase C Real Power", "M_AC_Power_SF", 0},
			{89, 1, "M_AC_Power_SF", TypeSF, "", "Real Power scale factor", "", 0},
			// ── Apparent Power ────────────────────────────────────────────────
			{90, 1, "M_AC_VA", TypeInt16, "VA", "Total Apparent Power", "M_AC_VA_SF", 0},
			{91, 1, "M_AC_VA_A", TypeInt16, "VA", "Phase A Apparent Power", "M_AC_VA_SF", 0},
			{92, 1, "M_AC_VA_B", TypeInt16, "VA", "Phase B Apparent Power", "M_AC_VA_SF", 0},
			{93, 1, "M_AC_VA_C", TypeInt16, "VA", "Phase C Apparent Power", "M_AC_VA_SF", 0},
			{94, 1, "M_AC_VA_SF", TypeSF, "", "Apparent Power scale factor", "", 0},
			// ── Reactive Power ────────────────────────────────────────────────
			{95, 1, "M_AC_VAR", TypeInt16, "VAR", "Total Reactive Power", "M_AC_VAR_SF", 0},
			{96, 1, "M_AC_VAR_A", TypeInt16, "VAR", "Phase A Reactive Power", "M_AC_VAR_SF", 0},
			{97, 1, "M_AC_VAR_B", TypeInt16, "VAR", "Phase B Reactive Power", "M_AC_VAR_SF", 0},
			{98, 1, "M_AC_VAR_C", TypeInt16, "VAR", "Phase C Reactive Power", "M_AC_VAR_SF", 0},
			{99, 1, "M_AC_VAR_SF", TypeSF, "", "Reactive Power scale factor", "", 0},
			// ── Power Factor ──────────────────────────────────────────────────
			{100, 1, "M_AC_PF", TypeInt16, "%", "Average Power Factor", "M_AC_PF_SF", 0},
			{101, 1, "M_AC_PF_A", TypeInt16, "%", "Phase A Power Factor", "M_AC_PF_SF", 0},
			{102, 1, "M_AC_PF_B", TypeInt16, "%", "Phase B Power Factor", "M_AC_PF_SF", 0},
			{103, 1, "M_AC_PF_C", TypeInt16, "%", "Phase C Power Factor", "M_AC_PF_SF", 0},
			{104, 1, "M_AC_PF_SF", TypeSF, "", "Power Factor scale factor", "", 0},
			// ── Real Energy ───────────────────────────────────────────────────
			{105, 2, "M_Exported", TypeAcc32, "Wh", "Total Exported Real Energy", "M_Energy_W_SF", 0},
			{107, 2, "M_Exported_A", TypeAcc32, "Wh", "Phase A Exported Real Energy", "M_Energy_W_SF", 0},
			{109, 2, "M_Exported_B", TypeAcc32, "Wh", "Phase B Exported Real Energy", "M_Energy_W_SF", 0},
			{111, 2, "M_Exported_C", TypeAcc32, "Wh", "Phase C Exported Real Energy", "M_Energy_W_SF", 0},
			{113, 2, "M_Imported", TypeAcc32, "Wh", "Total Imported Real Energy", "M_Energy_W_SF", 0},
			{115, 2, "M_Imported_A", TypeAcc32, "Wh", "Phase A Imported Real Energy", "M_Energy_W_SF", 0},
			{117, 2, "M_Imported_B", TypeAcc32, "Wh", "Phase B Imported Real Energy", "M_Energy_W_SF", 0},
			{119, 2, "M_Imported_C", TypeAcc32, "Wh", "Phase C Imported Real Energy", "M_Energy_W_SF", 0},
			{121, 1, "M_Energy_W_SF", TypeSF, "", "Real Energy scale factor", "", 0},
			// ── Apparent Energy ───────────────────────────────────────────────
			{122, 2, "M_Exported_VA", TypeAcc32, "VAh", "Total Exported Apparent Energy", "M_Energy_VA_SF", 0},
			{130, 2, "M_Imported_VA", TypeAcc32, "VAh", "Total Imported Apparent Energy", "M_Energy_VA_SF", 0},
			{138, 1, "M_Energy_VA_SF", TypeSF, "", "Apparent Energy scale factor", "", 0},
			// ── Reactive Energy (summary quadrants only) ──────────────────────
			{139, 2, "M_Import_VARh_Q1", TypeAcc32, "VARh", "Q1 Total Imported Reactive Energy", "M_Energy_VAR_SF", 0},
			{147, 2, "M_Import_VARh_Q2", TypeAcc32, "VARh", "Q2 Total Imported Reactive Energy", "M_Energy_VAR_SF", 0},
			{155, 2, "M_Export_VARh_Q3", TypeAcc32, "VARh", "Q3 Total Exported Reactive Energy", "M_Energy_VAR_SF", 0},
			{163, 2, "M_Export_VARh_Q4", TypeAcc32, "VARh", "Q4 Total Exported Reactive Energy", "M_Energy_VAR_SF", 0},
			{171, 1, "M_Energy_VAR_SF", TypeSF, "", "Reactive Energy scale factor", "", 0},
			// ── Events ────────────────────────────────────────────────────────
			{172, 2, "M_Events", TypeBitfield, "", "Meter event flags", "", 0},
		},
	}
}

// meterBase maps slot number to 0-based Modbus address.
func meterBase(slot int) uint16 {
	switch slot {
	case 2:
		return 40295
	case 3:
		return 40469
	default:
		return 40121
	}
}

// batteryBaseAddress returns the 0-based Modbus base address for a battery slot.
// Slot 1 → 0xE100, Slot 2 → 0xE200.
func batteryBaseAddress(slot int) uint16 {
	if slot == 2 {
		return 0xE200
	}
	return 0xE100
}

// BatteryInfoModel returns the static nameplate block for a SolarEdge battery.
// This is NOT a SunSpec standard model — it is SolarEdge proprietary,
// documented by community reverse-engineering (nmakel/solaredge_modbus).
//
// This model covers only the contiguous nameplate registers (0x00–0x4B),
// which include manufacturer, model, serial, firmware, and rated capacity.
// It is kept separate from BatteryDataModel to avoid reading across the
// 32-register gap at 0x4C–0x6B that causes timeouts on some firmware versions.
func BatteryInfoModel(slot int) *Model {
	return &Model{
		Name:        "BatteryInfo",
		BaseAddress: batteryBaseAddress(slot),
		WordOrder:   LowWordFirst,
		Fields: []RegisterDef{
			// ── Nameplate / identification ────────────────────────────────
			{0x00, 16, "B_Manufacturer", TypeString, "", "Battery manufacturer", "", 32},
			{0x10, 16, "B_Model", TypeString, "", "Battery model", "", 32},
			{0x20, 16, "B_Version", TypeString, "", "Battery firmware version", "", 32},
			{0x30, 16, "B_SerialNumber", TypeString, "", "Battery serial number", "", 32},
			{0x40, 1, "B_DeviceAddress", TypeUint16, "", "Modbus device address", "", 0},
			{0x41, 1, "B_SunSpec_DID", TypeUint16, "", "SunSpec device identifier", "", 0},
			// ── Rated capacity ────────────────────────────────────────────
			{0x42, 2, "B_RatedEnergy", TypeFloat32, "Wh", "Rated energy capacity", "", 0},
			{0x44, 2, "B_MaxChargePower", TypeFloat32, "W", "Maximum charge continuous power", "", 0},
			{0x46, 2, "B_MaxDischargePower", TypeFloat32, "W", "Maximum discharge continuous power", "", 0},
			{0x48, 2, "B_MaxChargePeakPower", TypeFloat32, "W", "Maximum charge peak power", "", 0},
			{0x4A, 2, "B_MaxDischargePeakPower", TypeFloat32, "W", "Maximum discharge peak power", "", 0},
		},
	}
}

// BatteryDataModel returns the live data block for a SolarEdge battery.
// This model covers the contiguous live data registers (0x6C–0x8D),
// which include temperature, voltage, current, power, energy counters,
// SoC, SoH, and status. It starts after the 32-register gap at 0x4C–0x6B.
func BatteryDataModel(slot int) *Model {
	// Live data starts at offset 0x6C from the battery base address.
	return &Model{
		Name:        "BatteryData",
		BaseAddress: batteryBaseAddress(slot) + 0x6C,
		WordOrder:   LowWordFirst,
		Fields: []RegisterDef{
			// ── Temperature ───────────────────────────────────────────────
			{0x00, 2, "B_TempAverage", TypeFloat32, "°C", "Average temperature", "", 0},
			{0x02, 2, "B_TempMax", TypeFloat32, "°C", "Maximum temperature", "", 0},
			// ── Live measurements ─────────────────────────────────────────
			{0x04, 2, "B_Voltage", TypeFloat32, "V", "Instantaneous voltage", "", 0},
			{0x06, 2, "B_Current", TypeFloat32, "A", "Instantaneous current (+ charge, - discharge)", "", 0},
			{0x08, 2, "B_Power", TypeFloat32, "W", "Instantaneous power (+ charge, - discharge)", "", 0},
			// ── Energy counters ───────────────────────────────────────────
			{0x0A, 4, "B_ExportedEnergy", TypeUint64, "Wh", "Lifetime exported (discharged) energy", "", 0},
			{0x0E, 4, "B_ImportedEnergy", TypeUint64, "Wh", "Lifetime imported (charged) energy", "", 0},
			// ── State ─────────────────────────────────────────────────────
			{0x12, 2, "B_EnergyMax", TypeFloat32, "Wh", "Maximum energy (at rated conditions)", "", 0},
			{0x14, 2, "B_EnergyAvailable", TypeFloat32, "Wh", "Available energy (current charge)", "", 0},
			{0x16, 2, "B_SoH", TypeFloat32, "%", "State of Health", "", 0},
			{0x18, 2, "B_SoC", TypeFloat32, "%", "State of Charge", "", 0},
			// Status and event fields are uint16 (1 register each), not uint32.
			// Reading them as uint32 merges two registers producing garbage values.
			{0x1A, 1, "B_Status", TypeUint16, "", "Battery status (0=Off 1=Standby 2=Init 3=Charging 4=Discharging 5=Fault 6=Idle)", "", 0},
			{0x1B, 1, "B_StatusInternal", TypeUint16, "", "Vendor battery status", "", 0},
			{0x1C, 1, "B_EventLog", TypeUint16, "", "Event log bitmask", "", 0},
			{0x1D, 1, "B_EventLogInternal", TypeUint16, "", "Vendor event log bitmask", "", 0},
		},
	}
}

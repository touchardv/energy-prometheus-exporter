package sunspec

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestDecodeInt16WithScaleFactor(t *testing.T) {
	// Simulate: I_AC_Power = 2071, I_AC_Power_SF = -2  → 20.71 W
	m := &Model{
		Name:        "TestInverter",
		BaseAddress: 0,
		Fields: []RegisterDef{
			{0, 1, "I_AC_Power", TypeInt16, "W", "AC Power", "I_AC_Power_SF", 0},
			{1, 1, "I_AC_Power_SF", TypeSF, "", "Scale factor", "", 0},
		},
	}
	raw := buildRaw(2071, 0xFFFE) // 0xFFFE is -2 in two's complement
	fields, err := Decode(m, raw)
	if err != nil {
		t.Fatal(err)
	}
	want := 20.71
	got := fields[0].Value
	if math.Abs(got-want) > 0.001 {
		t.Errorf("I_AC_Power: want %.3f, got %.3f", want, got)
	}
}

func TestDecodeAcc32WithScaleFactor(t *testing.T) {
	// Simulate: M_Exported = 12345678 Wh, M_Energy_W_SF = 0 → 12345678 Wh
	m := &Model{
		Name:        "TestMeter",
		BaseAddress: 0,
		Fields: []RegisterDef{
			{0, 2, "M_Exported", TypeAcc32, "Wh", "Exported Energy", "M_Energy_W_SF", 0},
			{2, 1, "M_Energy_W_SF", TypeSF, "", "Scale factor", "", 0},
		},
	}
	raw := buildRaw(0x00BC, 0x614E, 0x0000) // 0x00BC614E = 12345678, SF=0
	fields, err := Decode(m, raw)
	if err != nil {
		t.Fatal(err)
	}
	want := 12345678.0
	got := fields[0].Value
	if math.Abs(got-want) > 1 {
		t.Errorf("M_Exported: want %.0f, got %.0f", want, got)
	}
}

func TestDecodeString(t *testing.T) {
	m := &Model{
		Name:        "TestCommon",
		BaseAddress: 0,
		Fields: []RegisterDef{
			{0, 8, "C_Manufacturer", TypeString, "", "Manufacturer", "", 16},
		},
	}
	// "SolarEdge" padded with zeros to 16 bytes (8 registers)
	b := make([]byte, 16)
	copy(b, "SolarEdge")
	fields, err := Decode(m, b)
	if err != nil {
		t.Fatal(err)
	}
	if fields[0].StrVal != "SolarEdge" {
		t.Errorf("C_Manufacturer: want %q, got %q", "SolarEdge", fields[0].StrVal)
	}
}

func TestDecodeFloat32(t *testing.T) {
	// IEEE 754: 285.6 = 0x438ECCCD
	m := &Model{
		Name:        "TestBattery",
		BaseAddress: 0,
		Fields: []RegisterDef{
			{0, 2, "B_SoC", TypeFloat32, "%", "Battery SoC", "", 0},
		},
	}
	bits := math.Float32bits(285.6)
	raw := buildRaw(uint16(bits>>16), uint16(bits&0xFFFF))
	fields, err := Decode(m, raw)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(fields[0].Value-285.6) > 0.01 {
		t.Errorf("B_SoC: want 285.6, got %f", fields[0].Value)
	}
}

func buildRaw(regs ...uint16) []byte {
	b := make([]byte, len(regs)*2)
	for i, r := range regs {
		binary.BigEndian.PutUint16(b[i*2:], r)
	}
	return b
}

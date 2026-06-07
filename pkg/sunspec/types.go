package sunspec

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// WordOrder defines the order of 16-bit words within a multi-register value.
// Standard SunSpec devices use HighWordFirst (big-endian word order).
// SolarEdge proprietary battery registers use LowWordFirst.
type WordOrder int

const (
	HighWordFirst WordOrder = iota // standard: register N = high word (default)
	LowWordFirst                   // SolarEdge battery: register N = low word
)

// DataType mirrors the types used in the SunSpec / SolarEdge register tables.
type DataType string

const (
	TypeUint16   DataType = "uint16"
	TypeInt16    DataType = "int16"
	TypeUint32   DataType = "uint32"
	TypeAcc32    DataType = "acc32" // uint32 accumulator, same wire format
	TypeFloat32  DataType = "float32"
	TypeString   DataType = "string"
	TypeUint64   DataType = "uint64"
	TypeBitfield DataType = "bitfield32"
	TypeSF       DataType = "sunssf" // scale factor — int16 wire, special semantics
)

// RegistersNeeded returns the number of 16-bit Modbus registers occupied by
// a value of this type with the given string length (only relevant for TypeString).
func (t DataType) RegistersNeeded(strLen int) uint16 {
	switch t {
	case TypeUint16, TypeInt16, TypeSF:
		return 1
	case TypeUint32, TypeAcc32, TypeFloat32:
		return 2
	case TypeUint64:
		return 4
	case TypeBitfield:
		return 2
	case TypeString:
		return uint16((strLen + 1) / 2)
	}
	return 1
}

// RegisterDef describes a single field in a SunSpec model block.
// It is the Go equivalent of one row in the PDF tables.
type RegisterDef struct {
	// Offset is the 0-based position within the model block (not the absolute
	// Modbus address). Absolute address = model.BaseAddress + Offset.
	Offset uint16

	// Size in 16-bit registers as stated in the spec.
	Size uint16

	// Name is the field identifier from the spec, e.g. "M_AC_Power".
	Name string

	// Type is the data type.
	Type DataType

	// Unit is the physical unit, e.g. "Watts", "Amps", "°C". Empty if N/A.
	Unit string

	// Description is the human-readable label from the spec.
	Description string

	// ScaleFactorRef is the Name of another RegisterDef in the same model
	// whose value is the scale factor for this field. Empty if not applicable.
	ScaleFactorRef string

	// StringLen is the number of ASCII characters for TypeString fields.
	// Ignored for all other types.
	StringLen int
}

// Model represents a complete SunSpec model block (e.g. "Meter", "Inverter",
// "CommonBlock"). It holds the field definitions and, at read time, the raw
// bytes fetched from the device.
type Model struct {
	Name        string
	BaseAddress uint16    // absolute 0-based Modbus address of the first register
	WordOrder   WordOrder // word order for multi-register values (default: HighWordFirst)
	Fields      []RegisterDef
}

// TotalRegisters returns the number of registers needed to cover all fields,
// calculated as max(field.Offset + field.Size) across all fields.
// This is correct even when fields have gaps between them (as in the SunSpec
// Common block where offsets are non-contiguous).
func (m *Model) TotalRegisters() uint16 {
	var total uint16
	for _, f := range m.Fields {
		if end := f.Offset + f.Size; end > total {
			total = end
		}
	}
	return total
}

// NOT_IMPLEMENTED sentinel values per SunSpec spec.
const (
	notImplUint16 = 0xFFFF
	notImplInt16  = int16(-32768) // 0x8000
	notImplUint32 = 0xFFFFFFFF
	notImplSF     = int16(-32768) // scale factor also uses 0x8000
)

// DecodedField holds a fully decoded register value.
type DecodedField struct {
	Def            RegisterDef
	Raw            []byte  // raw bytes from the wire
	Value          float64 // numeric value (after scale factor, if any)
	StrVal         string  // populated for TypeString
	IsStr          bool
	NotImplemented bool // true when the device returned a NOT_IMPLEMENTED sentinel
}

// String returns a human-readable representation.
func (d DecodedField) String() string {
	if d.IsStr {
		return fmt.Sprintf("%-30s = %q  (%s)", d.Def.Name, d.StrVal, d.Def.Description)
	}
	if d.NotImplemented {
		return fmt.Sprintf("%-30s = %10s  %-10s  (%s)", d.Def.Name, "N/A", d.Def.Unit, d.Def.Description)
	}
	unit := d.Def.Unit
	if unit == "" {
		unit = "-"
	}
	return fmt.Sprintf("%-30s = %10.3f  %-10s  (%s)", d.Def.Name, d.Value, unit, d.Def.Description)
}

// isNotImplemented checks whether the raw register bytes represent a SunSpec
// NOT_IMPLEMENTED sentinel value for the given type.
func isNotImplemented(t DataType, chunk []byte) bool {
	switch t {
	case TypeUint16:
		return binary.BigEndian.Uint16(chunk[0:2]) == notImplUint16
	case TypeInt16:
		return int16(binary.BigEndian.Uint16(chunk[0:2])) == notImplInt16
	case TypeSF:
		return int16(binary.BigEndian.Uint16(chunk[0:2])) == notImplSF
	case TypeUint32, TypeAcc32:
		return binary.BigEndian.Uint32(chunk[0:4]) == notImplUint32
	case TypeUint64:
		return binary.BigEndian.Uint64(chunk[0:8]) == 0xFFFFFFFFFFFFFFFF
	case TypeBitfield:
		return binary.BigEndian.Uint32(chunk[0:4]) == notImplUint32
	}
	return false
}

// Decode takes the raw byte slice read from the device (one byte per register
// byte, big-endian) and decodes every field in the model.
// It returns a slice of DecodedField in the same order as model.Fields.
// wordSwap swaps the two 16-bit words in a 4-byte slice when LowWordFirst is used.
func wordSwap4(b []byte) []byte {
	return []byte{b[2], b[3], b[0], b[1]}
}

// wordSwap8 swaps word pairs for 8-byte (uint64) values under LowWordFirst.
func wordSwap8(b []byte) []byte {
	return []byte{b[6], b[7], b[4], b[5], b[2], b[3], b[0], b[1]}
}

func Decode(m *Model, raw []byte) ([]DecodedField, error) {
	// First pass: decode all raw numeric values (without scale factors).
	rawValues := make(map[string]float64, len(m.Fields))
	results := make([]DecodedField, len(m.Fields))

	for i, f := range m.Fields {
		byteOffset := int(f.Offset) * 2
		byteLen := int(f.Size) * 2

		if byteOffset+byteLen > len(raw) {
			return nil, fmt.Errorf("field %q: byte range [%d:%d] exceeds raw data length %d",
				f.Name, byteOffset, byteOffset+byteLen, len(raw))
		}

		chunk := raw[byteOffset : byteOffset+byteLen]
		df := DecodedField{Def: f, Raw: chunk}

		// Check for NOT_IMPLEMENTED sentinel before decoding.
		if f.Type != TypeString && isNotImplemented(f.Type, chunk) {
			df.NotImplemented = true
			rawValues[f.Name] = 0 // use 0 as safe default for SF lookups
			results[i] = df
			continue
		}

		switch f.Type {
		case TypeUint16:
			df.Value = float64(binary.BigEndian.Uint16(chunk[0:2]))

		case TypeInt16, TypeSF:
			df.Value = float64(int16(binary.BigEndian.Uint16(chunk[0:2])))

		case TypeUint32, TypeAcc32:
			w := chunk[0:4]
			if m.WordOrder == LowWordFirst {
				w = wordSwap4(w)
			}
			df.Value = float64(binary.BigEndian.Uint32(w))

		case TypeFloat32:
			w := chunk[0:4]
			if m.WordOrder == LowWordFirst {
				w = wordSwap4(w)
			}
			df.Value = float64(math.Float32frombits(binary.BigEndian.Uint32(w)))

		case TypeUint64:
			w := chunk[0:8]
			if m.WordOrder == LowWordFirst {
				w = wordSwap8(w)
			}
			df.Value = float64(binary.BigEndian.Uint64(w))

		case TypeBitfield:
			df.Value = float64(binary.BigEndian.Uint32(chunk[0:4]))

		case TypeString:
			b := make([]byte, 0, len(chunk))
			for _, c := range chunk {
				if c == 0 {
					break
				}
				b = append(b, c)
			}
			df.StrVal = strings.TrimRight(string(b), " ")
			df.IsStr = true
		}

		rawValues[f.Name] = df.Value
		results[i] = df
	}

	// Second pass: apply scale factors.
	for i, df := range results {
		if df.IsStr || df.Def.ScaleFactorRef == "" || df.NotImplemented {
			continue
		}
		sf, ok := rawValues[df.Def.ScaleFactorRef]
		if !ok {
			return nil, fmt.Errorf("field %q: scale factor ref %q not found in model",
				df.Def.Name, df.Def.ScaleFactorRef)
		}
		results[i].Value = df.Value * math.Pow10(int(sf))
	}

	return results, nil
}

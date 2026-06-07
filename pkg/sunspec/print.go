package sunspec

import (
	"fmt"
	"io"
	"strings"
)

// Print writes a formatted table of decoded fields to w.
func Print(w io.Writer, m *Model, fields []DecodedField) {
	sep := strings.Repeat("─", 80)
	fmt.Fprintf(w, "%s\n", sep)
	fmt.Fprintf(w, "  Model: %s  (base address: 0x%04X / %d)\n", m.Name, m.BaseAddress, m.BaseAddress)
	fmt.Fprintf(w, "%s\n", sep)
	fmt.Fprintf(w, "  %-30s  %-14s  %-10s  %s\n", "Name", "Value", "Unit", "Description")
	fmt.Fprintf(w, "%s\n", sep)
	for _, f := range fields {
		if f.Def.Type == TypeSF {
			// Print scale factors in a subdued way
			fmt.Fprintf(w, "  %-30s  %+14.0f  %-10s  %s\n",
				f.Def.Name, f.Value, "(SF)", f.Def.Description)
			continue
		}
		switch {
		case f.IsStr:
			fmt.Fprintf(w, "  %-30s  %-26s  %s\n",
				f.Def.Name, fmt.Sprintf("%q", f.StrVal), f.Def.Description)
		case f.NotImplemented:
			fmt.Fprintf(w, "  %-30s  %14s  %-10s  %s\n",
				f.Def.Name, "N/A", f.Def.Unit, f.Def.Description)
		default:
			unit := f.Def.Unit
			if unit == "" {
				unit = "-"
			}
			fmt.Fprintf(w, "  %-30s  %14.3f  %-10s  %s\n",
				f.Def.Name, f.Value, unit, f.Def.Description)
		}
	}
	fmt.Fprintf(w, "%s\n", sep)
}

// FieldMap returns a name→DecodedField lookup map for convenient access.
func FieldMap(fields []DecodedField) map[string]DecodedField {
	m := make(map[string]DecodedField, len(fields))
	for _, f := range fields {
		m[f.Def.Name] = f
	}
	return m
}

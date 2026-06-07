package sunspec

import (
	"fmt"

	"github.com/simonvetter/modbus"
)

// maxRegistersPerRead is the maximum registers per request.
// simonvetter/modbus caps this at 123 (the Modbus spec allows 125).
const maxRegistersPerRead = 123

// Reader reads a complete Model from a Modbus device, splitting into multiple
// requests when the model spans more than 125 registers.
type Reader struct {
	client *modbus.ModbusClient
}

func NewReader(client *modbus.ModbusClient) *Reader {
	return &Reader{client: client}
}

// Read fetches all registers for the model (chunking if necessary) and returns
// the decoded fields.
func (r *Reader) Read(m *Model) ([]DecodedField, error) {
	total := m.TotalRegisters()
	if total == 0 {
		return nil, fmt.Errorf("model %q has no fields", m.Name)
	}

	raw, err := r.readChunked(m.BaseAddress, total)
	if err != nil {
		return nil, fmt.Errorf("reading model %q at 0x%04X (%d regs): %w",
			m.Name, m.BaseAddress, total, err)
	}

	return Decode(m, raw)
}

// readChunked reads `count` registers starting at `addr`, issuing multiple
// Modbus requests of at most maxRegistersPerRead registers each and
// concatenating the results into a single byte slice.
func (r *Reader) readChunked(addr, count uint16) ([]byte, error) {
	raw := make([]byte, 0, int(count)*2)
	remaining := count
	offset := uint16(0)

	for remaining > 0 {
		chunkSize := remaining
		if chunkSize > maxRegistersPerRead {
			chunkSize = maxRegistersPerRead
		}

		// ReadRawBytes quantity is in bytes, not registers — pass chunkSize*2.
		chunk, err := r.client.ReadRawBytes(addr+offset, chunkSize*2, modbus.HOLDING_REGISTER)
		if err != nil {
			return nil, fmt.Errorf("chunk at offset %d (%d regs): %w", offset, chunkSize, err)
		}
		raw = append(raw, chunk...)

		offset += chunkSize
		remaining -= chunkSize
	}

	return raw, nil
}

// ReadField fetches and decodes a single named field from a model.
// The entire model block is read (chunked if needed) to resolve scale factor
// references correctly.
func (r *Reader) ReadField(m *Model, fieldName string) (*DecodedField, error) {
	fields, err := r.Read(m)
	if err != nil {
		return nil, err
	}
	for i, f := range fields {
		if f.Def.Name == fieldName {
			return &fields[i], nil
		}
	}
	return nil, fmt.Errorf("field %q not found in model %q", fieldName, m.Name)
}

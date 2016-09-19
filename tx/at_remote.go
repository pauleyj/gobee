package tx

import (
	"errors"
	"bytes"
)

const api_id_at_remote byte = 0x17

type AT_REMOTE struct {
	ID byte
	Addr64 uint64
	Addr16 uint16
	Options byte
	Command []byte
	Parameter []byte
}

func (f *AT_REMOTE) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer
	b.WriteByte(api_id_at_remote)
	b.WriteByte(f.ID)
	b.WriteByte(byte((f.Addr64 >> 56) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 48) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 40) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 32) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 24) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 16) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 8) & 0xFF))
	b.WriteByte(byte(f.Addr64 & 0xFF))
	b.WriteByte(byte((f.Addr16 >> 8) & 0xFF))
	b.WriteByte(byte(f.Addr16 & 0xFF))
	b.WriteByte(f.Options)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}
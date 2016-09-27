package tx

import (
	"bytes"
	"errors"
)

const atAPIID byte = 0x08

var _ Frame = (*AT)(nil)

// AT transmit frame
type AT struct {
	ID        byte
	Command   []byte
	Parameter []byte
}

// Bytes turn AT frame into bytes
func (f *AT) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer
	b.WriteByte(atAPIID)
	b.WriteByte(f.ID)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

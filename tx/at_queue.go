package tx

import (
	"bytes"
	"errors"
)

const atQueueAPIID byte = 0x09

var _ Frame = (*ATQueue)(nil)

// ATQueue AT queue transmit frame
type ATQueue struct {
	ID        byte
	Command   []byte
	Parameter []byte
}

// Bytes turn AT frame into bytes
func (f *ATQueue) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer

	b.WriteByte(atQueueAPIID)
	b.WriteByte(f.ID)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

package tx

import (
	"bytes"
	"errors"
)

const api_id_at_queue byte = 0x09

var _ TxFrame = (*AT_QUEUE)(nil)

type AT_QUEUE struct {
	ID        byte
	Command   []byte
	Parameter []byte
}

func (f *AT_QUEUE) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer
	b.WriteByte(api_id_at_queue)
	b.WriteByte(f.ID)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

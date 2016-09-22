package tx

import (
	"bytes"
	"errors"
)

const api_id_at byte = 0x08

var _ TxFrame = (*AT)(nil)

type AT struct {
	ID        byte
	Command   []byte
	Parameter []byte
}

func (f *AT) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer
	b.WriteByte(api_id_at)
	b.WriteByte(f.ID)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

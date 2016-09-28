package tx

import (
	"bytes"
	"errors"
)

const atRemoteAPIID byte = 0x17

var _ Frame = (*ATRemote)(nil)

// ATRemote AT remote transmit frame
type ATRemote struct {
	ID        byte
	Addr64    uint64
	Addr16    uint16
	Options   byte
	Command   []byte
	Parameter []byte
}

// Bytes turn ATRemote frame into bytes
func (f *ATRemote) Bytes() ([]byte, error) {
	if len(f.Command) != 2 {
		return nil, errors.New("Invalid AT command")
	}

	var b bytes.Buffer

	b.WriteByte(atRemoteAPIID)
	b.WriteByte(f.ID)
	b.Write(uint64ToBytes(f.Addr64))
	b.Write(uint16ToBytes(f.Addr16))
	b.WriteByte(f.Options)
	b.Write(f.Command)

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

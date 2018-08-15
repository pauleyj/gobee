package tx

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const atRemoteAPIID byte = 0x17

func NewATRemote(options ...func(interface{})) *ATRemote {
	f := &ATRemote{Addr64:0xFFFF, Addr16:0xFFFE, Cmd:[2]byte{'N','I'}}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// ATRemote AT remote transmit frame
type ATRemote struct {
	FrameID byte
	Addr64 uint64
	Addr16 uint16
	Options byte
	Cmd [2]byte
	Parameter []byte
}

func (f *ATRemote) SetFrameID(id byte) {
	f.FrameID = id
}

func (f *ATRemote) SetAddr64(addr uint64) {
	f.Addr64 = addr
}

func (f *ATRemote) SetAddr16(addr uint16) {
	f.Addr16 = addr
}

func (f *ATRemote) SetOptions(options byte) {
	f.Options = options
}

func (f *ATRemote) SetCommand(cmd [2]byte) {
	copy(f.Cmd[:], cmd[:])
}

func (f *ATRemote) SetParameter(parameter []byte) {
	f.Parameter = make([]byte, len(parameter))
	copy(f.Parameter, parameter)
}

// Bytes turn ATRemote frame into bytes
func (f *ATRemote) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(atRemoteAPIID)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.FrameID)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint64ToBytes(f.Addr64))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint16ToBytes(f.Addr16))
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.Options)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(f.Cmd[:])
	if err != nil {
		return nil, err
	}

	if f.Parameter != nil && len(f.Parameter) > 0 {
		_, err = b.Write(f.Parameter)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

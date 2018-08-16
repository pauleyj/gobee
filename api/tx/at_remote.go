package tx

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const atRemoteAPIID byte = 0x17

func NewATRemote(options ...func(interface{})) *ATRemote {
	f := &ATRemote{Addr64:0xFFFF, Addr16:0xFFFE, Cmd:[2]byte{'N','I'}}

	optionsRunner(f, options...)

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

	b.WriteByte(atRemoteAPIID)
	b.WriteByte(f.FrameID)
	b.Write(util.Uint64ToBytes(f.Addr64))
	b.Write(util.Uint16ToBytes(f.Addr16))
	b.WriteByte(f.Options)
	b.Write(f.Cmd[:])

	if f.Parameter != nil && len(f.Parameter) > 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

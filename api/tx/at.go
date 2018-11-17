package tx

import (
	"bytes"
)

const atAPIID byte = 0x08

func NewAT(options ...func(interface{})) *AT {
	f := &AT{Cmd: NI}

	optionsRunner(f, options...)

	return f
}

// AT transmit frame
type AT struct {
	FrameID   byte
	Cmd       [2]byte
	Parameter []byte
}

// SetFrameID satisfy FrameIDSetter interface
func (f *AT) SetFrameID(id byte) {
	f.FrameID = id
}

// SetCommand satisfy CommandSetter interface
func (f *AT) SetCommand(cmd [2]byte) {
	copy(f.Cmd[:], cmd[:])
}

// SetParameter satisfy ParameterSetter interface
func (f *AT) SetParameter(parameter []byte) {
	f.Parameter = make([]byte, len(parameter))
	copy(f.Parameter, parameter)
}

// Bytes turn AT frame into bytes, satisfy Frame interface
func (f *AT) Bytes() ([]byte, error) {
	var b bytes.Buffer

	b.WriteByte(atAPIID)
	b.WriteByte(f.FrameID)
	b.Write(f.Cmd[:])

	if f.Parameter != nil && len(f.Parameter) > 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

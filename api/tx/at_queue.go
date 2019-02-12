package tx

import (
	"bytes"
)

const atQueueAPIID byte = 0x09

func NewATQueue(options ...func(interface{})) *ATQueue {
	f := &ATQueue{}

	optionsRunner(f, options...)

	return f
}

// ATQueue queue transmit frame
type ATQueue struct {
	FrameID   byte
	Cmd       [2]byte
	Parameter []byte
}

// SetFrameID satisfy FrameIDSetter interface
func (f *ATQueue) SetFrameID(id byte) {
	f.FrameID = id
}

// SetCommand satisfy CommandSetter interface
func (f *ATQueue) SetCommand(cmd [2]byte) {
	copy(f.Cmd[:], cmd[:])
}

// SetParameter satisfy ParameterSetter interface
func (f *ATQueue) SetParameter(parameter []byte) {
	f.Parameter = make([]byte, len(parameter))
	copy(f.Parameter, parameter)
}

// Bytes turn ATQueue frame into bytes, satisfy Frame interface
func (f *ATQueue) Bytes() ([]byte, error) {
	var b bytes.Buffer

	b.WriteByte(atQueueAPIID)
	b.WriteByte(f.FrameID)
	b.Write(f.Cmd[:])

	if f.Parameter != nil {
		b.Write(f.Parameter)
	}

	return b.Bytes(), nil
}

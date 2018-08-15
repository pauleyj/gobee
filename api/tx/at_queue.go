package tx

import (
	"bytes"
)

const atQueueAPIID byte = 0x09

func NewATQueue(options ...func(interface{})) *ATQueue {
	f := &ATQueue{}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// ATQueue queue transmit frame
type ATQueue struct {
	FrameID   byte
	Cmd       [2]byte
	Parameter []byte
}

func (f *ATQueue) SetFrameID(id byte) {
	f.FrameID = id
}

func (f *ATQueue) SetCommand(cmd [2]byte) {
	copy(f.Cmd[:], cmd[:])
}

func (f *ATQueue) SetParameter(parameter []byte) {
	f.Parameter = make([]byte, len(parameter))
	copy(f.Parameter, parameter)
}

// Bytes turn ATQueue frame into bytes
func (f *ATQueue) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(atQueueAPIID)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.FrameID)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(f.Cmd[:])
	if err != nil {
		return nil, err
	}

	if f.Parameter != nil {
		_, err = b.Write(f.Parameter)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

package at

import (
	"bytes"
)

const atAPIID byte = 0x08

func FrameID(id byte) func(*AT) {
	return func(f *AT) {
		f.FrameID = id
	}
}

func Command(cmd [2]byte) func(*AT) {
	return func(f *AT) {
		copy(f.Cmd[:], cmd[:])
	}
}

func Parameter(parameter []byte) func(*AT) {
	return func(f *AT) {
		if parameter == nil || len(parameter) == 0 {
			return
		}

		f.Parameter = make([]byte, len(parameter))
		copy(f.Parameter, parameter)
	}
}

func NewAT(options ...func(*AT)) *AT {
	f := &AT{Cmd:[2]byte{'N','I'}}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// AT transmit frame
type AT struct {
	FrameID   byte
	Cmd       [2]byte
	Parameter []byte
}

// Bytes turn AT frame into bytes
func (f *AT) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(atAPIID)
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

	if f.Parameter != nil && len(f.Parameter) > 0 {
		_, err = b.Write(f.Parameter)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

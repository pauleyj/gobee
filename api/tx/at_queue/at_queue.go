package at_queue

import (
	"bytes"
)

const atQueueAPIID byte = 0x09

func FrameID(id byte) func(*ATQueue) {
	return func(f *ATQueue) {
		f.FrameID = id
	}
}

func Command(cmd [2]byte) func(*ATQueue) {
	return func(f *ATQueue) {
		copy(f.Cmd[:], cmd[:])
		for i, b := range cmd {
			f.Cmd[i] = b
		}
	}
}

func Parameter(parameter []byte) func(*ATQueue) {
	return func(f *ATQueue) {
		if parameter == nil || len(parameter) == 0 {
			return
		}

		f.Parameter = make([]byte, len(parameter))
		copy(f.Parameter, parameter)
	}
}

func NewATQueue(options ...func(*ATQueue)) *ATQueue {
	f := &ATQueue{}

	for _, option := range options {
		option(f)
	}

	return f
}

// ATQueue queue transmit frame
type ATQueue struct {
	buffer bytes.Buffer
	
	FrameID   byte
	Cmd       [2]byte
	Parameter []byte
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

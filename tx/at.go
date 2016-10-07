package tx

import (
	"bytes"
)

const atAPIID byte = 0x08

var _ Frame = (*AT)(nil)

// NewATBuilder builder of AT frames
func NewATBuilder() *atID {
	return &atID{}
}

type atID struct {
	buffer bytes.Buffer
}

func (b *atID) ID(id byte) *atCommand {
	b.buffer.WriteByte(atAPIID)
	b.buffer.WriteByte(id)
	return &atCommand{buffer: b.buffer}
}

type atCommand struct {
	buffer bytes.Buffer
}

func (b *atCommand) Command(command [2]byte) *atParameter {
	b.buffer.Write(command[:])
	return &atParameter{buffer: b.buffer}
}

type atParameter struct {
	buffer bytes.Buffer
}

func (b *atParameter) Parameter(parameter *byte) *atBuilder {
	if parameter != nil {
		b.buffer.WriteByte(*parameter)
	}
	return &atBuilder{buffer: b.buffer}
}

type atBuilder struct {
	buffer bytes.Buffer
}

func (b *atBuilder) Build() *AT {
	return &AT{buffer: b.buffer}
}

// AT transmit frame
type AT struct {
	buffer bytes.Buffer
}

// Bytes turn AT frame into bytes
func (f *AT) Bytes() ([]byte, error) {
	return f.buffer.Bytes(), nil
}

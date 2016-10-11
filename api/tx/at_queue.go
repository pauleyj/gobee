package tx

import (
	"bytes"
)

const atQueueAPIID byte = 0x09

var _ Frame = (*ATQueue)(nil)

// NewATQueueBuilder builder of ATQueue frames
func NewATQueueBuilder() *atQueueID {
	return &atQueueID{}
}

type atQueueID struct {
	buffer bytes.Buffer
}

func (b *atQueueID) ID(id byte) *atQueueCommand {
	b.buffer.WriteByte(atQueueAPIID)
	b.buffer.WriteByte(id)
	return &atQueueCommand{buffer: b.buffer}
}

type atQueueCommand struct {
	buffer bytes.Buffer
}

func (b *atQueueCommand) Command(command [2]byte) *atQueueParameter {
	b.buffer.Write(command[:])
	return &atQueueParameter{buffer: b.buffer}
}

type atQueueParameter struct {
	buffer bytes.Buffer
}

func (b *atQueueParameter) Parameter(parameter *byte) *atQueueBuilder {
	if parameter != nil {
		b.buffer.WriteByte(*parameter)
	}
	return &atQueueBuilder{buffer: b.buffer}
}

type atQueueBuilder struct {
	buffer bytes.Buffer
}

func (b *atQueueBuilder) Build() *ATQueue {
	return &ATQueue{buffer: b.buffer}
}

// ATQueue AT queue transmit frame
type ATQueue struct {
	buffer bytes.Buffer
}

// Bytes turn AT frame into bytes
func (f *ATQueue) Bytes() ([]byte, error) {
	return f.buffer.Bytes(), nil
}

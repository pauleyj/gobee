package tx

import (
	"bytes"
)

const atRemoteAPIID byte = 0x17

var _ Frame = (*ATRemote)(nil)

// NewATRemoteBuilder builder of ATRemote frames
func NewATRemoteBuilder() *atRemoteID {
	return &atRemoteID{}
}

type atRemoteID struct {
	buffer bytes.Buffer
}

func (b *atRemoteID) ID(id byte) *atRemoteAddr64 {
	b.buffer.WriteByte(atRemoteAPIID)
	b.buffer.WriteByte(id)
	return &atRemoteAddr64{buffer: b.buffer}
}

type atRemoteAddr64 struct {
	buffer bytes.Buffer
}

func (b *atRemoteAddr64) Addr64(addr uint64) *atRemoteAddr16 {
	b.buffer.Write(uint64ToBytes(addr))
	return &atRemoteAddr16{buffer: b.buffer}
}

type atRemoteAddr16 struct {
	buffer bytes.Buffer
}

func (b *atRemoteAddr16) Addr16(addr uint16) *atRemoteOptions {
	b.buffer.Write(uint16ToBytes(addr))
	return &atRemoteOptions{buffer: b.buffer}
}

type atRemoteOptions struct {
	buffer bytes.Buffer
}

func (b *atRemoteOptions) Options(options byte) *atRemoteCommand {
	b.buffer.WriteByte(options)
	return &atRemoteCommand{buffer: b.buffer}
}

type atRemoteCommand struct {
	buffer bytes.Buffer
}

func (b *atRemoteCommand) Command(command [2]byte) *atRemoteParameter {
	b.buffer.Write(command[:])
	return &atRemoteParameter{buffer: b.buffer}
}

type atRemoteParameter struct {
	buffer bytes.Buffer
}

func (b *atRemoteParameter) Parameter(parameter *byte) *atRemoteBuilder {
	if parameter != nil {
		b.buffer.WriteByte(*parameter)
	}
	return &atRemoteBuilder{buffer: b.buffer}
}

type atRemoteBuilder struct {
	buffer bytes.Buffer
}

func (b *atRemoteBuilder) Build() *ATRemote {
	return &ATRemote{buffer: b.buffer}
}

// ATRemote AT remote transmit frame
type ATRemote struct {
	buffer bytes.Buffer
}

// Bytes turn ATRemote frame into bytes
func (f *ATRemote) Bytes() ([]byte, error) {
	return f.buffer.Bytes(), nil
}

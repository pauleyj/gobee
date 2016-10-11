package tx

import (
	"bytes"
)

const zbAPIID byte = 0x10

var _ Frame = (*ZB)(nil)

// NewZBBuilder builder of ZB frames
func NewZBBuilder() *zbID {
	return &zbID{}
}

type zbID struct {
	buffer bytes.Buffer
}

func (b *zbID) ID(id byte) *zbAddr64 {
	b.buffer.WriteByte(zbAPIID)
	b.buffer.WriteByte(id)
	return &zbAddr64{buffer: b.buffer}
}

type zbAddr64 struct {
	buffer bytes.Buffer
}

func (b *zbAddr64) Addr64(addr uint64) *zbAddr16 {
	b.buffer.Write(uint64ToBytes(addr))
	return &zbAddr16{buffer: b.buffer}
}

type zbAddr16 struct {
	buffer bytes.Buffer
}

func (b *zbAddr16) Addr16(addr uint16) *zbBroadcastRadius {
	b.buffer.Write(uint16ToBytes(addr))
	return &zbBroadcastRadius{buffer: b.buffer}
}

type zbBroadcastRadius struct {
	buffer bytes.Buffer
}

func (b *zbBroadcastRadius) BroadcastRadius(broadcastRadius byte) *zbOptions {
	b.buffer.WriteByte(broadcastRadius)
	return &zbOptions{buffer: b.buffer}
}

type zbOptions struct {
	buffer bytes.Buffer
}

func (b *zbOptions) Options(options byte) *zbData {
	b.buffer.WriteByte(options)
	return &zbData{buffer: b.buffer}
}

type zbData struct {
	buffer bytes.Buffer
}

func (b *zbData) Data(data []byte) *zbBuilder {
	if len(data) > 0 {
		b.buffer.Write(data)
	}
	return &zbBuilder{buffer: b.buffer}
}

type zbBuilder struct {
	buffer bytes.Buffer
}

func (b *zbBuilder) Build() *ZB {
	return &ZB{buffer: b.buffer}
}

// ZB transmit frame
type ZB struct {
	buffer bytes.Buffer
}

// Bytes turn ATRemote frame into bytes
func (f *ZB) Bytes() ([]byte, error) {
	return f.buffer.Bytes(), nil
}

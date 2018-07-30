package tx

import (
	"bytes"
)

const zbExplicitAPIID byte = 0x11

var _ Frame = (*ZBExplicit)(nil)

// NewZBExplicitBuilder builder of ZBExplicit frames
func NewZBExplicitBuilder() *zbExplicitID {
	return &zbExplicitID{}
}

type zbExplicitID struct {
	buffer bytes.Buffer
}

func (b *zbExplicitID) ID(id byte) *zbExplicitAddr64 {
	b.buffer.WriteByte(zbExplicitAPIID)
	b.buffer.WriteByte(id)
	return &zbExplicitAddr64{buffer: b.buffer}
}

type zbExplicitAddr64 struct {
	buffer bytes.Buffer
}

func (b *zbExplicitAddr64) Addr64(addr uint64) *zbExplicitAddr16 {
	//b.buffer.Write(uint64ToBytes(addr))
	return &zbExplicitAddr16{buffer: b.buffer}
}

type zbExplicitAddr16 struct {
	buffer bytes.Buffer
}

func (b *zbExplicitAddr16) Addr16(addr uint16) *zbExplicitSrcEP {
	//b.buffer.Write(uint16ToBytes(addr))
	return &zbExplicitSrcEP{buffer: b.buffer}
}

type zbExplicitSrcEP struct {
	buffer bytes.Buffer
}

func (b *zbExplicitSrcEP) SrcEP(endpoint byte) *zbExplicitDstEP {
	b.buffer.WriteByte(endpoint)
	return &zbExplicitDstEP{buffer: b.buffer}
}

type zbExplicitDstEP struct {
	buffer bytes.Buffer
}

func (b *zbExplicitDstEP) DstEP(endpoint byte) *zbExplicitClusterID {
	b.buffer.WriteByte(endpoint)
	return &zbExplicitClusterID{buffer: b.buffer}
}

type zbExplicitClusterID struct {
	buffer bytes.Buffer
}

func (b *zbExplicitClusterID) ClusterID(id uint16) *zbExplicitProfileID {
	//b.buffer.Write(uint16ToBytes(id))
	return &zbExplicitProfileID{buffer: b.buffer}
}

type zbExplicitProfileID struct {
	buffer bytes.Buffer
}

func (b *zbExplicitProfileID) ProfileID(id uint16) *zbExplicitBroadcastRadius {
	//b.buffer.Write(uint16ToBytes(id))
	return &zbExplicitBroadcastRadius{buffer: b.buffer}
}

type zbExplicitBroadcastRadius struct {
	buffer bytes.Buffer
}

func (b *zbExplicitBroadcastRadius) BroadcastRadius(broadcastRadius byte) *zbExplicitOptions {
	b.buffer.WriteByte(broadcastRadius)
	return &zbExplicitOptions{buffer: b.buffer}
}

type zbExplicitOptions struct {
	buffer bytes.Buffer
}

func (b *zbExplicitOptions) Options(options byte) *zbExplicitData {
	b.buffer.WriteByte(options)
	return &zbExplicitData{buffer: b.buffer}
}

type zbExplicitData struct {
	buffer bytes.Buffer
}

func (b *zbExplicitData) Data(data []byte) *zbExplicitBuilder {
	if len(data) > 0 {
		b.buffer.Write(data)
	}
	return &zbExplicitBuilder{buffer: b.buffer}
}

type zbExplicitBuilder struct {
	buffer bytes.Buffer
}

func (b *zbExplicitBuilder) Build() *ZBExplicit {
	return &ZBExplicit{buffer: b.buffer}
}

// ZBExplicit transmit frame
type ZBExplicit struct {
	buffer bytes.Buffer
}

// Bytes turn ATRemote frame into bytes
func (f *ZBExplicit) Bytes() ([]byte, error) {
	return f.buffer.Bytes(), nil
}

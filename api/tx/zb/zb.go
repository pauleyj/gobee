package zb

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const zbAPIID byte = 0x10

func FrameID(id byte) func(*ZB) {
	return func(f *ZB) {
		f.FrameID = id
	}
}

func Addr64(addr64 uint64) func(*ZB) {
	return func(f *ZB) {
		f.Addr64 = addr64
	}
}

func Addr16(addr16 uint16) func(*ZB) {
	return func(f *ZB) {
		f.Addr16 = addr16
	}
}

func BroadcastRadius(hops byte) func(*ZB) {
	return func(f *ZB) {
		f.BroadcastRadius = hops
	}
}

func Options(options byte) func(*ZB) {
	return func(f *ZB) {
		f.Options = options
	}
}

func Data(data []byte) func(*ZB) {
	return func(f *ZB) {
		if data == nil || len(data) == 0 {
			return
		}

		f.Data = make([]byte, len(data))
		copy(f.Data, data)
	}
}

// ZB transmit frame
type ZB struct {
	FrameID         byte
	Addr64          uint64
	Addr16          uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func NewZB(options ...func(*ZB)) *ZB {
	f := &ZB{Addr64:0xFFFF, Addr16:0xFFFE}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// Bytes turn ZB frame into bytes
func (f *ZB) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(zbAPIID)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.FrameID)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint64ToBytes(f.Addr64))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint16ToBytes(f.Addr16))
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.BroadcastRadius)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.Options)
	if err != nil {
		return nil, err
	}

	if f.Data != nil && len(f.Data) > 0 {
		_, err = b.Write(f.Data)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

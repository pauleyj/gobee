package zb_explicit

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const zbExplicitAPIID byte = 0x11

func FrameID(id byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.FrameID = id
	}
}

func Addr64(addr64 uint64) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.Addr64 = addr64
	}
}

func Addr16(addr16 uint16) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.Addr16 = addr16
	}
}

func SrcEP(src byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.SrcEP = src
	}
}
func DstEP(dst byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.DstEP = dst
	}
}
func ClusterID(id byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.ClusterID = id
	}
}
func ProfileID(id byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.ProfileID = id
	}
}

func BroadcastRadius(hops byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.BroadcastRadius = hops
	}
}

func Options(options byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		f.Options = options
	}
}

func Data(data []byte) func(*ZBExplicit) {
	return func(f *ZBExplicit) {
		if data == nil || len(data) == 0 {
			return
		}

		f.Data = make([]byte, len(data))
		copy(f.Data, data)
	}
}

// ZBExplicit transmit frame
type ZBExplicit struct {
	FrameID         byte
	Addr64          uint64
	Addr16          uint16
	SrcEP           byte
	DstEP           byte
	ClusterID       byte
	ProfileID       byte
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func NewZBExplicit(options ...func(*ZBExplicit)) *ZBExplicit {
	f := &ZBExplicit{Addr64:0xFFFF, Addr16:0xFFFE}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// Bytes turn ATRemote frame into bytes
func (f *ZBExplicit) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(zbExplicitAPIID)
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

	err = b.WriteByte(f.SrcEP)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.DstEP)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.ClusterID)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.ProfileID)
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

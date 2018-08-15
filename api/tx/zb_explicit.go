package tx

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const zbExplicitAPIID byte = 0x11

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

func NewZBExplicit(options ...func(interface{})) *ZBExplicit {
	f := &ZBExplicit{Addr64:0xFFFF, Addr16:0xFFFE}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

func (f *ZBExplicit) SetFrameID(id byte) {
	f.FrameID = id
}

func (f *ZBExplicit) SetAddr64(addr uint64) {
	f.Addr64 = addr
}

func (f *ZBExplicit) SetAddr16(addr uint16) {
	f.Addr16 = addr
}

func (f *ZBExplicit) SetSrcEP(src byte) {
	f.SrcEP = src
}

func (f *ZBExplicit) SetDstEP(dst byte) {
	f.DstEP = dst
}

func (f *ZBExplicit) SetClusterID(id byte) {
	f.ClusterID = id
}

func (f *ZBExplicit) SetProfileID(id byte) {
	f.ProfileID = id
}

func (f *ZBExplicit) SetBroadcastRadius(hops byte) {
	f.BroadcastRadius = hops
}

func (f *ZBExplicit) SetOptions(options byte) {
	f.Options = options
}

func (f *ZBExplicit) SetData(data []byte) {
	f.Data = make([]byte, len(data))
	copy(f.Data, data)
}

// Bytes turn frame into bytes
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

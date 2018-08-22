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
	ClusterID       uint16
	ProfileID       uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func NewZBExplicit(options ...func(interface{})) *ZBExplicit {
	f := &ZBExplicit{Addr64:0xFFFF, Addr16:0xFFFE}

	optionsRunner(f, options...)

	return f
}

// SetFrameID satisfy FrameIDSetter interface
func (f *ZBExplicit) SetFrameID(id byte) {
	f.FrameID = id
}

// SetAddr64 satisfy Addr64Setter interface
func (f *ZBExplicit) SetAddr64(addr uint64) {
	f.Addr64 = addr
}

// SetAddr16 satisfy Addr16Setter interface
func (f *ZBExplicit) SetAddr16(addr uint16) {
	f.Addr16 = addr
}

// SetSrcEP satisfy SrcEPSetter interface
func (f *ZBExplicit) SetSrcEP(src byte) {
	f.SrcEP = src
}

// SetDstEP satisfy DstEPSetter interface
func (f *ZBExplicit) SetDstEP(dst byte) {
	f.DstEP = dst
}

// SetClusterID satisfy ClusterIDSetter interface
func (f *ZBExplicit) SetClusterID(id uint16) {
	f.ClusterID = id
}

// SetProfileID satisfy ProfileIDSetter interface
func (f *ZBExplicit) SetProfileID(id uint16) {
	f.ProfileID = id
}

// SetBroadcastRadius satisfy BroadcastRadiusSetter interface
func (f *ZBExplicit) SetBroadcastRadius(hops byte) {
	f.BroadcastRadius = hops
}

// SetOptions satisfy OptionsSetter interface
func (f *ZBExplicit) SetOptions(options byte) {
	f.Options = options
}

// SetData satisfy DataSetter interface
func (f *ZBExplicit) SetData(data []byte) {
	f.Data = make([]byte, len(data))
	copy(f.Data, data)
}

// Bytes turn frame into bytes, satosfy Frame interface
func (f *ZBExplicit) Bytes() ([]byte, error) {
	var b bytes.Buffer

	b.WriteByte(zbExplicitAPIID)
	b.WriteByte(f.FrameID)
	b.Write(util.Uint64ToBytes(f.Addr64))
	b.Write(util.Uint16ToBytes(f.Addr16))
	b.WriteByte(f.SrcEP)
	b.WriteByte(f.DstEP)
	b.Write(util.Uint16ToBytes(f.ClusterID))
	b.Write(util.Uint16ToBytes(f.ProfileID))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if f.Data != nil && len(f.Data) > 0 {
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}

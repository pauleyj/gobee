package tx

import (
	"bytes"
)

const zbExplicitAPIID byte = 0x11

var _ Frame = (*ZBExplicit)(nil)

// ZBExplicit transmit frame
type ZBExplicit struct {
	ID              byte
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

// Bytes turn ATRemote frame into bytes
func (f *ZBExplicit) Bytes() ([]byte, error) {
	b := new(bytes.Buffer)

	b.WriteByte(zbExplicitAPIID)
	b.WriteByte(f.ID)
	b.Write(uint64ToBytes(f.Addr64))
	b.Write(uint16ToBytes(f.Addr16))
	b.WriteByte(f.SrcEP)
	b.WriteByte(f.DstEP)
	b.Write(uint16ToBytes(f.ClusterID))
	b.Write(uint16ToBytes(f.ProfileID))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) > 0 {
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}

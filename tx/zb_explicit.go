package tx

import (
	"bytes"
)

const api_id_zb_explicit byte = 0x11

var _ TxFrame = (*ZB_EXPLICIT)(nil)

type ZB_EXPLICIT struct {
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

func (f *ZB_EXPLICIT) Bytes() ([]byte, error) {
	b := new(bytes.Buffer)

	b.WriteByte(api_id_zb_explicit)
	b.WriteByte(f.ID)
	b.WriteByte(byte((f.Addr64 >> 56) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 48) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 40) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 32) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 24) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 16) & 0xFF))
	b.WriteByte(byte((f.Addr64 >> 8) & 0xFF))
	b.WriteByte(byte(f.Addr64 & 0xFF))
	b.WriteByte(byte((f.Addr16 >> 8) & 0xFF))
	b.WriteByte(byte(f.Addr16 & 0xFF))
	b.WriteByte(f.SrcEP)
	b.WriteByte(f.DstEP)
	b.WriteByte(byte((f.ClusterID >> 8) & 0xFF))
	b.WriteByte(byte(f.ClusterID & 0xFF))
	b.WriteByte(byte((f.ProfileID >> 8) & 0xFF))
	b.WriteByte(byte(f.ProfileID & 0xFF))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) > 0{
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}
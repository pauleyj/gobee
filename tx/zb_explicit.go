package tx

import (
	"bytes"
	"errors"
)

const api_id_zb_explicit byte = 0x11

type ZB_EXPLICIT struct {
	ID              byte
	Addr64          uint64
	Addr16          uint16
	SrcEP           byte
	DstEP           byte
	ClusterID       []byte
	ProfileID       []byte
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func (f *ZB_EXPLICIT) Bytes() ([]byte, error) {
	if len(f.ClusterID) != 2 {
		return nil, errors.New("Invalid Cluster ID")
	}
	if len(f.ProfileID) != 2 {
		return nil, errors.New("Invalid Profile ID")
	}

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
	b.Write(f.ClusterID)
	b.Write(f.ProfileID)
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) > 0{
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}
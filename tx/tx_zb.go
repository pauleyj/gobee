package tx

import "bytes"

const api_id_tx_zb byte = 0x10

type TX_ZB struct {
	ID              byte
	DestAddr64      uint64
	DestAddr16      uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func (f *TX_ZB) Bytes() []byte {
	b := new(bytes.Buffer)

	b.WriteByte(api_id_tx_zb)
	b.WriteByte(f.ID)
	b.WriteByte(byte((f.DestAddr64 >> 56) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 48) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 40) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 32) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 24) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 16) & 0xFF))
	b.WriteByte(byte((f.DestAddr64 >> 8) & 0xFF))
	b.WriteByte(byte(f.DestAddr64 & 0xFF))
	b.WriteByte(byte((f.DestAddr16 >> 8) & 0xFF))
	b.WriteByte(byte(f.DestAddr16 & 0xFF))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) != 0 {
		b.Write(f.Data)
	}

	return b.Bytes()
}
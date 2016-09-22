package tx

import "bytes"

const api_id_tx_zb byte = 0x10

var _ TxFrame = (*ZB)(nil)

type ZB struct {
	ID              byte
	Addr64          uint64
	Addr16          uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func (f *ZB) Bytes() ([]byte, error) {
	b := new(bytes.Buffer)

	b.WriteByte(api_id_tx_zb)
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
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) != 0 {
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}
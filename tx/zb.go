package tx

import "bytes"

const zbAPIID byte = 0x10

var _ Frame = (*ZB)(nil)

// ZB transmit frame
type ZB struct {
	ID              byte
	Addr64          uint64
	Addr16          uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

// Bytes turn ATRemote frame into bytes
func (f *ZB) Bytes() ([]byte, error) {
	var b bytes.Buffer

	b.WriteByte(zbAPIID)
	b.WriteByte(f.ID)
	b.Write(uint64ToBytes(f.Addr64))
	b.Write(uint16ToBytes(f.Addr16))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if len(f.Data) != 0 {
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}

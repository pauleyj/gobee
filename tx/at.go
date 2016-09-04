package tx

import "bytes"

const api_id_at byte = 0x08

type AT struct {
	ID        byte
	Command   [2]byte
	Parameter []byte
}

func (f *AT) Bytes() []byte {
	var b bytes.Buffer

	b.WriteByte(api_id_at)
	b.WriteByte(f.ID)
	b.Write(f.Command[:])

	if len(f.Parameter) != 0 {
		b.Write(f.Parameter)
	}

	return b.Bytes()
}

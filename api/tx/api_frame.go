package tx

import (
	"bytes"
	"github.com/pauleyj/gobee/api"
)

// APIFrame defines an API frame structure
type APIFrame struct {
	Mode api.EscapeMode
}

// Bytes transforms a data frame into an API frame slice
func (f *APIFrame) Bytes(frame Frame) ([]byte, error) {
	p, err := frame.Bytes()
	if err != nil {
		return p, err
	}

	var b bytes.Buffer
	b.WriteByte(api.FrameDelimiter)
	b.Write(f.encode(uint16ToBytes(uint16(len(p)))...))
	b.Write(f.encode(p...))
	c := checksum(p)
	b.Write(f.encode(c))

	return b.Bytes(), nil
}

func (f *APIFrame) encode(p ...byte) []byte {
	if f.Mode == api.EscapeModeInactive {
		return slicify(p...)
	}
	return escape(p...)
}

func escape(p ...byte) []byte {
	var b bytes.Buffer
	for _, c := range p {
		if api.ShouldEscape(c) {
			b.WriteByte(api.ESC)
			b.WriteByte(api.Escape(c))
		} else {
			b.WriteByte(c)
		}
	}
	return b.Bytes()
}

func slicify(p ...byte) []byte {
	var b bytes.Buffer
	for _, c := range p {
		b.WriteByte(c)
	}
	return b.Bytes()
}

func checksum(p []byte) byte {
	var chksum byte
	for _, c := range p {
		chksum += c
	}

	return api.ValidChecksum - chksum
}

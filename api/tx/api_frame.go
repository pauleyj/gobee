package tx

import (
	"bytes"

	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/tx/util"
)

// New constructs a TX API frame
func New(options ...func(interface{})) *APIFrame {
	f := &APIFrame{}

	optionsRunner(f, options...)

	return f
}

// APIFrame defines an API frame structure
type APIFrame struct {
	Mode api.EscapeMode
}

// SetAPIEscapeMode satisfy APIEscapeModeSetter interface
func (f *APIFrame) SetAPIEscapeMode(mode api.EscapeMode) {
	f.Mode = mode
}

// Bytes transforms a data frame into an API frame slice
func (f *APIFrame) Bytes(frame Frame) ([]byte, error) {
	p, err := frame.Bytes()
	if err != nil {
		return p, err
	}

	var b bytes.Buffer
	b.WriteByte(api.FrameDelimiter)
	b.Write(f.encode(util.Uint16ToBytes(uint16(len(p)))))
	b.Write(f.encode(p))
	b.Write(f.encode([]byte{checksum(p)}))

	return b.Bytes(), nil
}

func (f *APIFrame) encode(p []byte) []byte {
	if f.Mode == api.EscapeModeInactive {
		return p
	}
	return escape(p)
}

func escape(p []byte) []byte {
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

func checksum(p []byte) byte {
	var chksum byte
	for _, c := range p {
		chksum += c
	}

	return api.ValidChecksum - chksum
}

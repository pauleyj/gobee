package rx

import "encoding/binary"

const (
	zbAPIID byte = 0x90

	zbAddr64Offset  = 0
	zbAddr16Offset  = 8
	zbOptionsOffset = 10
	zbDataOffset    = 11
)

var _ Frame = (*ZB)(nil)

// ZB rx frame
type ZB struct {
	buffer []byte
}

func newZB() Frame {
	return &ZB{
		buffer: make([]byte, 0),
	}
}

// RX frame data
func (f *ZB) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

// Addr64 64-bit address of sender
func (f *ZB) Addr64() uint64 {
	return binary.BigEndian.Uint64(f.buffer[zbAddr64Offset : zbAddr64Offset+addr64Length])
}

// Addr16 16-bit address of sender
func (f *ZB) Addr16() uint16 {
	return binary.BigEndian.Uint16(f.buffer[zbAddr16Offset : zbAddr16Offset+addr16Length])
}

// Options frame options
func (f *ZB) Options() byte {
	return f.buffer[zbOptionsOffset]
}

// Data frame data
func (f *ZB) Data() []byte {
	if len(f.buffer) == zbDataOffset {
		return nil
	}

	return f.buffer[zbDataOffset:]
}

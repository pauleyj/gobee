package rx

import "encoding/binary"

const (
	atRemoteAPIID byte = 0x97

	atRemoteIDOffset      = 0
	atRemoteAddr64Offset  = 1
	atRemoteAddr16Offset  = 9
	atRemoteCommandOffset = 11
	atRemoteCommandLength = 2
	atRemoteStatusOffset  = 13
	atRemoteDataOffset    = 14
)

var _ Frame = (*ATRemote)(nil)

// ATRemote rx frame
type ATRemote struct {
	buffer []byte
}

func newATRemote() Frame {
	return &ATRemote{
		buffer: make([]byte, 0),
	}
}

// RX frame data
func (f *ATRemote) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

// ID frame ID
func (f *ATRemote) ID() byte {
	return f.buffer[atRemoteIDOffset]
}

// Addr64 remote 64-bit address
func (f *ATRemote) Addr64() uint64 {
	return binary.BigEndian.Uint64(f.buffer[atRemoteAddr64Offset : atRemoteAddr64Offset+addr64Length])
}

// Addr16 remote 16-bit address
func (f *ATRemote) Addr16() uint16 {
	return binary.BigEndian.Uint16(f.buffer[atRemoteAddr16Offset : atRemoteAddr16Offset+addr16Length])
}

// Command remote AT command
func (f *ATRemote) Command() []byte {
	return f.buffer[atRemoteCommandOffset : atRemoteCommandOffset+atRemoteCommandLength]
}

// Status remote AT command status
func (f *ATRemote) Status() byte {
	return f.buffer[atRemoteStatusOffset]
}

// Data remote AT command data
func (f *ATRemote) Data() []byte {
	if len(f.buffer) == atRemoteDataOffset {
		return nil
	}

	return f.buffer[atRemoteDataOffset:]
}

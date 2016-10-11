package rx

const (
	atAPIID byte = 0x88

	atIDOffset      = 0
	atCommandOffset = 1
	atCommandLength = 2
	atStatusOffset  = 3
	atDataOffset    = 4
)

var _ Frame = (*AT)(nil)

// AT rx frame
type AT struct {
	buffer []byte
}

func newAT() Frame {
	return &AT{
		buffer: make([]byte, 0),
	}
}

// RX frame data
func (f *AT) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

// ID frame ID
func (f *AT) ID() byte {
	return f.buffer[atIDOffset]
}

// Command AT command
func (f *AT) Command() []byte {
	return f.buffer[atCommandOffset : atCommandOffset+atCommandLength]
}

// Status AT command status
func (f *AT) Status() byte {
	return f.buffer[atStatusOffset]
}

// Data AT command data
func (f *AT) Data() []byte {
	if len(f.buffer) == atDataOffset {
		return nil
	}

	return f.buffer[atDataOffset:]
}

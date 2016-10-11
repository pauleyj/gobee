package rx

import "encoding/binary"

const (
	txStatusAPIID byte = 0x8B

	txStatusFrameIDOffset         = 0
	txStatusAddr16Offset          = 1
	txStatusRetryCountOffset      = 3
	txStatusDeliveryStatusOffset  = 4
	txStatusDiscoveryStatusOffset = 5
)

var _ Frame = (*TXStatus)(nil)

// TXStatus rx frame
type TXStatus struct {
	buffer []byte
}

func newTXStatus() Frame {
	return &TXStatus{
		buffer: make([]byte, 0),
	}
}

// RX frame data
func (f *TXStatus) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

// ID frame ID of TX frame this status is associated with
func (f *TXStatus) ID() byte {
	return f.buffer[txStatusFrameIDOffset]
}

// Addr16 16-bit address of XBee this status message is associated with
func (f *TXStatus) Addr16() uint16 {
	return binary.BigEndian.Uint16(f.buffer[txStatusAddr16Offset : txStatusAddr16Offset+addr16Length])
}

// Retries number of retries
func (f *TXStatus) Retries() byte {
	return f.buffer[txStatusRetryCountOffset]
}

// Delivery delivery status of the TX
func (f *TXStatus) Delivery() byte {
	return f.buffer[txStatusDeliveryStatusOffset]
}

// Discovery discovery status
func (f *TXStatus) Discovery() byte {
	return f.buffer[txStatusDiscoveryStatusOffset]
}

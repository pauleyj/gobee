package rx

const txStatusAPIID byte = 0x8B

const (
	txStatusFrameID         = rxFrameState(iota)
	txStatusAddr16          = rxFrameState(iota)
	txStatusRetryCount      = rxFrameState(iota)
	txStatusDeliveryStatus  = rxFrameState(iota)
	txStatusDiscoveryStatus = rxFrameState(iota)
)

var _ Frame = (*TXStatus)(nil)

// TXStatus rx frame
type TXStatus struct {
	state     rxFrameState
	index     byte
	ID        byte
	Addr16    uint16
	Retries   byte
	Delivery  byte
	Discovery byte
}

func newTXStatus() Frame {
	return &TXStatus{
		state: txStatusFrameID,
	}
}

// RX frame data
func (f *TXStatus) RX(b byte) error {
	var err error

	switch f.state {
	case txStatusFrameID:
		err = f.stateFrameID(b)
	case txStatusAddr16:
		err = f.stateAddr16(b)
	case txStatusRetryCount:
		err = f.stateRetries(b)
	case txStatusDeliveryStatus:
		err = f.stateDeliveryStatus(b)
	case txStatusDiscoveryStatus:
		err = f.stateDiscoveryStatus(b)
	}

	return err

}

func (f *TXStatus) stateFrameID(b byte) error {
	f.ID = b
	f.state = txStatusAddr16

	return nil
}

func (f *TXStatus) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = txStatusRetryCount
	}

	return nil

}

func (f *TXStatus) stateRetries(b byte) error {
	f.Retries = b
	f.state = txStatusDeliveryStatus

	return nil
}

func (f *TXStatus) stateDeliveryStatus(b byte) error {
	f.Delivery = b
	f.state = txStatusDiscoveryStatus

	return nil
}

func (f *TXStatus) stateDiscoveryStatus(b byte) error {
	f.Discovery = b

	return nil
}

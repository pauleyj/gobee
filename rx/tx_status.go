package rx

const XBEE_API_TX_STATUS byte = 0x8B

const (
	tx_status_frame_id = rx_frame_state(iota)
	tx_status_addr16 = rx_frame_state(iota)
	tx_status_retry_count = rx_frame_state(iota)
	tx_status_delivery_status = rx_frame_state(iota)
	tx_status_discovery_status = rx_frame_state(iota)
)

type TX_STATUS struct {
	state rx_frame_state
	index byte
	ID byte
	Addr16 uint16
	Retries byte
	Delivery byte
	Discovery byte
}

func newTXStatus() RxFrame {
	return &TX_STATUS{
		state: tx_status_frame_id,
	}
}

func (f *TX_STATUS) RX(b byte) error {
	var err error

	switch f.state {
	case tx_status_frame_id:
		err = f.stateFrameID(b)
	case tx_status_addr16:
		err = f.stateAddr16(b)
	case tx_status_retry_count:
		err = f.stateRetries(b)
	case tx_status_delivery_status:
		err = f.stateDeliveryStatus(b)
	case tx_status_discovery_status:
		err = f.stateDiscoveryStatus(b)
	}

	return err

}

func (f *TX_STATUS) stateFrameID(b byte) error {
	f.ID = b
	f.state = tx_status_addr16

	return nil
}

func (f *TX_STATUS) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = tx_status_retry_count
	}

	return nil

}

func (f *TX_STATUS) stateRetries(b byte) error {
	f.Retries = b
	f.state = tx_status_delivery_status

	return nil
}

func (f *TX_STATUS) stateDeliveryStatus(b byte) error {
	f.Delivery = b
	f.state = tx_status_discovery_status

	return nil
}

func (f *TX_STATUS) stateDiscoveryStatus(b byte) error {
	f.Discovery = b

	return nil
}
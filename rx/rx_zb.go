package rx

const XBEE_API_ID_RX_ZB byte = 0x90

const (
	rx_zb_addr64 = rx_frame_state(iota)
	rx_zb_addr16 = rx_frame_state(iota)
	rx_zb_options = rx_frame_state(iota)
	rx_zb_data = rx_frame_state(iota)
)

type RX_ZB struct {
	state rx_frame_state
	index byte
	Addr64 uint64
	Addr16 uint16
	Options byte
	Data []byte
}

func newRX_ZB() RxFrame {
	return &RX_ZB{
		state: rx_zb_addr64,
	}
}

func (f *RX_ZB) RX(b byte) error {
	var err error

	switch f.state {
	case rx_zb_addr64:
		err = stateAddr64(f, b)
	case rx_zb_addr16:
		err = stateAddr16(f, b)
	case rx_zb_options:
		err = stateOptions(f, b)
	case rx_zb_data:
		err = stateData(f, b)
	}

	return err
}

func stateAddr64(f *RX_ZB, b byte) error  {
	f.Addr64 += uint64(b << (56 - (8 * f.index)))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = rx_zb_addr16
	}

	return nil
}

func stateAddr16(f *RX_ZB, b byte) error {
	f.Addr16 += uint16(b << (8 - (8 * f.index)))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = rx_zb_options
	}

	return nil

}

func stateOptions(f *RX_ZB, b byte) error {
	f.Options = b
	f.state = rx_zb_data

	return nil
}

func stateData(f *RX_ZB, b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}

	f.Data = append(f.Data, b)

	return nil
}
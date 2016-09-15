package rx

const XBEE_API_ID_RX_ZB byte = 0x90

const (
	rx_zb_addr64 = rx_frame_state(iota)
	rx_zb_addr16 = rx_frame_state(iota)
	rx_zb_options = rx_frame_state(iota)
	rx_zb_data = rx_frame_state(iota)
)

type ZB struct {
	state rx_frame_state
	index byte
	Addr64 uint64
	Addr16 uint16
	Options byte
	Data []byte
}

func newZB() RxFrame {
	return &ZB{
		state: rx_zb_addr64,
	}
}

func (f *ZB) RX(b byte) error {
	var err error

	switch f.state {
	case rx_zb_addr64:
		err = f.stateAddr64(b)
	case rx_zb_addr16:
		err = f.stateAddr16(b)
	case rx_zb_options:
		err = f.stateOptions(b)
	case rx_zb_data:
		err = f.stateData(b)
	}

	return err
}

func (f *ZB) stateAddr64(b byte) error  {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = rx_zb_addr16
	}

	return nil
}

func (f *ZB) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = rx_zb_options
	}

	return nil

}

func (f *ZB) stateOptions(b byte) error {
	f.Options = b
	f.state = rx_zb_data

	return nil
}

func (f *ZB) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}

	f.Data = append(f.Data, b)

	return nil
}
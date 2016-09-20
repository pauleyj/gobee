package rx

const XBEE_API_ID_RX_ZB_EXPLICIT byte = 0x91

const (
	rx_zbe_addr64 = rx_frame_state(iota)
	rx_zbe_addr16 = rx_frame_state(iota)
	rx_zbe_srcep = rx_frame_state(iota)
	rx_zbe_dstep = rx_frame_state(iota)
	rx_zbe_cid = rx_frame_state(iota)
	rx_zbe_pid = rx_frame_state(iota)
	rx_zbe_options = rx_frame_state(iota)
	rx_zbe_data = rx_frame_state(iota)
)

type ZB_EXPLICIT struct {
	state     rx_frame_state
	index     byte
	Addr64    uint64
	Addr16    uint16
	SrcEP     byte
	DstEP     byte
	ClusterID uint16
	ProfileID uint16
	Options   byte
	Data      []byte
}

func newZB_EXPLICIT() RxFrame {
	return &ZB_EXPLICIT{
		state: rx_zbe_addr64,
	}
}

func (f *ZB_EXPLICIT) RX(b byte) error {
	var err error

	switch f.state {
	case rx_zbe_addr64:
		err = f.stateAddr64(b)
	case rx_zbe_addr16:
		err = f.stateAddr16(b)
	case rx_zbe_srcep:
		err = f.stateSrcEP(b)
	case rx_zbe_dstep:
		err = f.stateDstEP(b)
	case rx_zbe_cid:
		err = f.stateCID(b)
	case rx_zbe_pid:
		err = f.statePID(b)
	case rx_zbe_options:
		err = f.stateOptions(b)
	case rx_zbe_data:
		err = f.stateData(b)
	}

	return err
}

func (f *ZB_EXPLICIT) stateAddr64(b byte) error {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = rx_zbe_addr16
	}

	return nil
}

func (f *ZB_EXPLICIT) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = rx_zbe_srcep
	}

	return nil

}

func (f *ZB_EXPLICIT) stateSrcEP(b byte) error {
	f.SrcEP = b
	f.state = rx_zbe_dstep

	return nil
}

func (f *ZB_EXPLICIT) stateDstEP(b byte) error {
	f.DstEP = b
	f.state = rx_zbe_cid

	return nil
}

func (f *ZB_EXPLICIT) stateCID(b byte) error {
	f.ClusterID += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = rx_zbe_pid
	}
	return nil
}

func (f *ZB_EXPLICIT) statePID(b byte) error {
	f.ProfileID += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = rx_zbe_options
	}
	return nil
}

func (f *ZB_EXPLICIT) stateOptions(b byte) error {
	f.Options = b
	f.state = rx_zbe_data

	return nil
}

func (f *ZB_EXPLICIT) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}

	f.Data = append(f.Data, b)

	return nil
}
package rx

const zbAPIID byte = 0x90

const (
	zbAddr64  = rxFrameState(iota)
	zbAddr16  = rxFrameState(iota)
	zbOptions = rxFrameState(iota)
	zbData    = rxFrameState(iota)
)

var _ Frame = (*ZB)(nil)

// ZB rx frame
type ZB struct {
	state   rxFrameState
	index   byte
	Addr64  uint64
	Addr16  uint16
	Options byte
	Data    []byte
}

func newZB() Frame {
	return &ZB{
		state: zbAddr64,
	}
}

// RX frame data
func (f *ZB) RX(b byte) error {
	var err error

	switch f.state {
	case zbAddr64:
		err = f.stateAddr64(b)
	case zbAddr16:
		err = f.stateAddr16(b)
	case zbOptions:
		err = f.stateOptions(b)
	case zbData:
		err = f.stateData(b)
	}

	return err
}

func (f *ZB) stateAddr64(b byte) error {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = zbAddr16
	}

	return nil
}

func (f *ZB) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = zbOptions
	}

	return nil

}

func (f *ZB) stateOptions(b byte) error {
	f.Options = b
	f.state = zbData

	return nil
}

func (f *ZB) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}

	f.Data = append(f.Data, b)

	return nil
}

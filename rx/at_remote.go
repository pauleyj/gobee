package rx

const atRemoteAPIID byte = 0x97

const (
	atRemoteID      = rxFrameState(iota)
	atRemoteAddr64  = rxFrameState(iota)
	atRemoteAddr16  = rxFrameState(iota)
	atRemoteCommand = rxFrameState(iota)
	atRemoteStatus  = rxFrameState(iota)
	atRemoteData    = rxFrameState(iota)
)

var _ Frame = (*ATRemote)(nil)

// ATRemote rx frame
type ATRemote struct {
	state   rxFrameState
	index   byte
	ID      byte
	Addr64  uint64
	Addr16  uint16
	Command [2]byte
	Status  byte
	Data    []byte
}

func newATRemote() Frame {
	return &ATRemote{
		state: atRemoteID,
	}
}

// RX frame data
func (f *ATRemote) RX(b byte) error {
	var err error

	switch f.state {
	case atRemoteID:
		err = f.stateID(b)
	case atRemoteAddr64:
		err = f.stateAddr64(b)
	case atRemoteAddr16:
		err = f.stateAddr16(b)
	case atRemoteCommand:
		err = f.stateCommand(b)
	case atRemoteStatus:
		err = f.stateStatus(b)
	case atRemoteData:
		err = f.stateData(b)
	}

	return err
}

func (f *ATRemote) stateID(b byte) error {
	f.ID = b
	f.state = atRemoteAddr64

	return nil
}

func (f *ATRemote) stateAddr64(b byte) error {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = atRemoteAddr16
	}

	return nil
}

func (f *ATRemote) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = atRemoteCommand
	}

	return nil

}

func (f *ATRemote) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == 2 {
		f.state = atRemoteStatus
	}
	return nil
}

func (f *ATRemote) stateStatus(b byte) error {
	f.Status = b
	f.state = atRemoteData

	return nil
}

func (f *ATRemote) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}

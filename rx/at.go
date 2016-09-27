package rx

const atAPIID byte = 0x88

const (
	atID      = rxFrameState(iota)
	atCommand = rxFrameState(iota)
	atStatus  = rxFrameState(iota)
	atData    = rxFrameState(iota)
)

var _ Frame = (*AT)(nil)

// AT rx frame
type AT struct {
	state   rxFrameState
	index   byte
	ID      byte
	Command [2]byte
	Status  byte
	Data    []byte
}

func newAT() Frame {
	return &AT{
		state: atID,
	}
}

// RX frame data
func (f *AT) RX(b byte) error {
	var err error

	switch f.state {
	case atID:
		err = f.stateID(b)
	case atCommand:
		err = f.stateCommand(b)
	case atStatus:
		err = f.stateStatus(b)
	case atData:
		err = f.stateData(b)
	}

	return err
}

func (f *AT) stateID(b byte) error {
	f.ID = b
	f.state = atCommand

	return nil
}

func (f *AT) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == 2 {
		f.state = atStatus
	}
	return nil
}

func (f *AT) stateStatus(b byte) error {
	f.Status = b
	f.state = atData

	return nil
}

func (f *AT) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}

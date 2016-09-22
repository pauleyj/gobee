package rx

const XBEE_API_ID_RX_AT byte = 0x88

const (
	at_id = rx_frame_state(iota)
	at_command = rx_frame_state(iota)
	at_status = rx_frame_state(iota)
	at_data = rx_frame_state(iota)
)

var _ RxFrame = (*AT)(nil)

type AT struct {
	state   rx_frame_state
	index   byte
	ID      byte
	Command [2]byte
	Status  byte
	Data    []byte
}

func newAT() RxFrame {
	return &AT{
		state:   at_id,
	}
}

func (f *AT) RX(b byte) error {
	var err error

	switch f.state {
	case at_id:
		err = f.stateId(b)
	case at_command:
		err = f.stateCommand(b)
	case at_status:
		err = f.stateStatus(b)
	case at_data:
		err = f.stateData(b)
	}

	return err
}

func (f *AT) stateId(b byte) error {
	f.ID = b
	f.state = at_command

	return nil
}

func (f *AT) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == 2 {
		f.state = at_status
	}
	return nil
}

func (f *AT) stateStatus(b byte) error {
	f.Status = b
	f.state = at_data

	return nil
}

func (f *AT) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}

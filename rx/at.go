package rx

const XBEE_API_ID_RX_AT byte = 0x88

const (
	at_state_frame_id = rx_frame_state(iota)
	at_state_frame_command_at = rx_frame_state(iota)
	at_state_frame_command_status = rx_frame_state(iota)
	at_state_frame_command_data = rx_frame_state(iota)
)

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
		state:   at_state_frame_id,
	}
}

func (f *AT) RX(b byte) error {
	var err error

	switch f.state {
	case at_state_frame_id:
		err = f.stateId(b)
	case at_state_frame_command_at:
		err = f.stateCommand(b)
	case at_state_frame_command_status:
		err = f.stateStatus(b)
	case at_state_frame_command_data:
		err = f.stateData(b)
	}

	return err
}

func (f *AT) stateId(b byte) error {
	f.ID = b
	f.state = at_state_frame_command_at

	return nil
}

func (f *AT) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == 2 {
		f.state = at_state_frame_command_status
	}
	return nil
}

func (f *AT) stateStatus(b byte) error {
	f.Status = b
	f.state = at_state_frame_command_data

	return nil
}

func (f *AT) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}

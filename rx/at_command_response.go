package rx

const XBEE_API_ID_AT_COMMAND_RESPONSE byte = 0x88
const at_command_length byte = 2

const (
	at_state_frame_id = rx_frame_state(iota)
	at_state_frame_command_at = rx_frame_state(iota)
	at_state_frame_command_status = rx_frame_state(iota)
	at_state_frame_command_data = rx_frame_state(iota)
)

type ATCommandResponse struct {
	state   rx_frame_state
	index   byte
	FrameId byte
	Command [2]byte
	Status  byte
	Data    []byte
}

func newATCommandResponse() RxFrame {
	return &ATCommandResponse{
		state:   at_state_frame_id,
	}
}

func (f *ATCommandResponse) RX(b byte) error {
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

func (f *ATCommandResponse) stateId(b byte) error {
	f.FrameId = b
	f.state = at_state_frame_command_at

	return nil
}

func (f *ATCommandResponse) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == at_command_length {
		f.state = at_state_frame_command_status
	}
	return nil
}

func (f *ATCommandResponse) stateStatus(b byte) error {
	f.Status = b
	f.state = at_state_frame_command_data

	return nil
}

func (f *ATCommandResponse) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}

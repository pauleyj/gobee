package rx

const XBEE_API_ID_RX_AT_REMOTE byte = 0x97

const (
	at_remote_id = rx_frame_state(iota)
	at_remote_addr64 = rx_frame_state(iota)
	at_remote_addr16 = rx_frame_state(iota)
	at_remote_command = rx_frame_state(iota)
	at_remote_status = rx_frame_state(iota)
	at_remote_data = rx_frame_state(iota)
)

type AT_REMOTE struct {
	state rx_frame_state
	index byte
	ID byte
	Addr64 uint64
	Addr16 uint16
	Command [2]byte
	Status byte
	Data []byte
}

func newAT_REMOTE() RxFrame {
	return &AT_REMOTE{
		state: at_remote_id,
	}
}

func (f *AT_REMOTE) RX(b byte) error {
	var err error

	switch f.state {
	case at_remote_id:
		err = f.stateID(b)
	case at_remote_addr64:
		err = f.stateAddr64(b)
	case at_remote_addr16:
		err = f.stateAddr16(b)
	case at_remote_command:
		err = f.stateCommand(b)
	case at_remote_status:
		err = f.stateStatus(b)
	case at_remote_data:
		err = f.stateData(b)
	}

	return err
}

func (f *AT_REMOTE) stateID(b byte) error {
	f.ID = b
	f.state = at_remote_addr64

	return nil
}
func (f *AT_REMOTE) stateAddr64(b byte) error {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = at_remote_addr16
	}

	return nil
}

func (f *AT_REMOTE) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = at_remote_command
	}

	return nil

}

func (f *AT_REMOTE) stateCommand(b byte) error {
	f.Command[f.index] = b
	f.index++
	if f.index == 2 {
		f.state = at_remote_status
	}
	return nil
}

func (f *AT_REMOTE) stateStatus(b byte) error {
	f.Status = b
	f.state = at_remote_data

	return nil
}

func (f *AT_REMOTE) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}
	f.Data = append(f.Data, b)

	return nil
}
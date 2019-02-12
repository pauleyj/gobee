package rx

import "fmt"

const (
	modemStatusAPIID byte = 0x8A
)

type ModemStatus struct {
	status byte
}

func newModemStatus() Frame {
	return &ModemStatus{}
}

func (f *ModemStatus) RX(b byte) error {
	f.status = b

	return nil
}

func (f *ModemStatus) Status() byte {
	return f.status
}

func (f *ModemStatus) String() string {
	switch f.status {
	case 0:
		return "hardware reset"
	case 1:
		return "watchdog timer reset"
	case 2:
		return "joined network (routers and end devices)"
	case 3:
		return "disassociated"
	case 6:
		return "coordinator started"
	case 7:
		return "network security key was updated"
	case 0x0d:
		return "voltage supply limit exceeded"
	case 0x11:
		return "modem configuration changed while join in progress"
	default:
		return fmt.Sprintf("unknown status (%#0.2x)", f.status)
	}
}

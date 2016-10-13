package gobee

import (
	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/rx"
	"github.com/pauleyj/gobee/api/tx"
)

// BroadcastAddr64 64-bit broadcast address
const BroadcastAddr64 uint64 = 0x000000000000FFFF

// BroadcastAddr16 16-bit broadcast address
const BroadcastAddr16 uint16 = 0xFFFE

// XBeeTransmitter used to transmit API frame bytes to serial communications port
type XBeeTransmitter interface {
	Transmit([]byte) (int, error)
}

// XBeeReceiver used to report frames received by the XBee
type XBeeReceiver interface {
	Receive(rx.Frame) error
}

// XBee all the things
type XBee struct {
	transmitter XBeeTransmitter
	receiver    XBeeReceiver
	apiMode     api.APIEscapeMode
	frame       *rx.APIFrame
}

// New constructor of XBee's
func New(transmitter XBeeTransmitter, receiver XBeeReceiver) *XBee {
	return &XBee{
		transmitter: transmitter,
		receiver:    receiver,
		apiMode:     api.EscapeModeInactive,
		frame:       &rx.APIFrame{Mode: api.EscapeModeInactive},
	}
}

// NewWithEscapeMode constructor of XBee's with a specific escape mode
func NewWithEscapeMode(transmitter XBeeTransmitter, receiver XBeeReceiver, mode api.APIEscapeMode) *XBee {
	return &XBee{
		transmitter: transmitter,
		receiver:    receiver,
		apiMode:     mode,
		frame:       &rx.APIFrame{Mode: mode},
	}
}

// RX bytes received from the serial communications port are sent here
func (x *XBee) RX(b byte) error {
	f, err := x.frame.RX(b)
	if err != nil {
		return err
	}

	if f != nil {
		x.receiver.Receive(f)
	}

	return nil
}

// TX transmit a frame to the XBee, forms an appropriate API frame for the frame being sent,
// uses the XBeeTransmitter to send the API frame bytes to the serial communications port
func (x *XBee) TX(frame tx.Frame) (int, error) {
	f := &tx.APIFrame{Mode: x.apiMode}
	p, err := f.Bytes(frame)
	if err != nil {
		return 0, err
	}

	return x.transmitter.Transmit(p)
}

// SetAPIMode sets the API mode so goobe knows to escape or not
func (x *XBee) SetAPIMode(mode api.APIEscapeMode) error {
	if mode != api.EscapeModeInactive && mode != api.EscapeModeActive {
		return api.ErrInvalidAPIEscapeMode
	}
	x.apiMode = mode
	return nil
}

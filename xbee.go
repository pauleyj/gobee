package gobee

import (
	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/rx"
	"github.com/pauleyj/gobee/api/tx"
)

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
	apiMode     api.EscapeMode
	frame       *rx.APIFrame
}

func APIEscapeMode(mode api.EscapeMode) func(*XBee) {
	return func(xbee *XBee) {
		xbee.apiMode = mode
	}
}

// New constructor of XBee's
func New(transmitter XBeeTransmitter, receiver XBeeReceiver, options ...func(xbee *XBee)) *XBee {
	xbee :=  &XBee{
		transmitter: transmitter,
		receiver:    receiver,
		frame:       &rx.APIFrame{},
	}

	if options == nil || len(options) == 0 {
		return xbee
	}

	for _, option := range options {
		option(xbee)
	}

	return xbee
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

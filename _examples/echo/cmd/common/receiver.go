package common

import "github.com/pauleyj/gobee/api/rx"

// NewReceiver constructs a new Receiver
func NewReceiver(rx chan<- rx.Frame) *Receiver {
	return &Receiver{
		rx: rx,
	}
}

// Receiver implements gobee.XBeeReceiver.
type Receiver struct {
	rx chan<- rx.Frame
}

// Receive satisfies gobee.XBeeReceiver interface.  This simple implementation
// simply puts the received frame onto rx channel.
func (r *Receiver) Receive(f rx.Frame) error {
	r.rx <- f

	return nil
}

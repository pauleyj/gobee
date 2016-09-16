package gobee_test

import (
	"testing"
	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
)

type Transmitter struct{
	t *testing.T
	expected []byte
}

func (xm *Transmitter) Transmit(p []byte) (n int, err error) {
	return len(p), nil
}

func (xm *Transmitter) SetExpectedWriteBytes(expected []byte) {
	xm.expected = expected
}

type Receiver struct {
	t *testing.T
}



func (r *Receiver) Receive(f rx.RxFrame) error {
	switch f.(type) {
	case *rx.AT:
		validateAT(r.t, f.(*rx.AT))

	}

	return nil
}

func validateAT(t *testing.T, f *rx.AT) {
	if f.FrameId != 1 {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02x", f.FrameId)
	}
}

func TestXBee_TX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	var at = &tx.AT{
		ID: 0x01,
		Parameter: []byte{0x00},
	}
	_, err := xbee.TX(at)
	if err == nil {
		t.Error("Expected error, but got none")
	}

	at = &tx.AT{
		ID: 0x01,
		Command: []byte{'N','I'},
		Parameter: []byte{0x00},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = &tx.AT{
		ID: 0x01,
		Command: []byte{'N','I'},
		Parameter: []byte{0x11},
	}
	xbee.SetApiMode(2)
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_AT_QUEUE(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	var at = &tx.AT_QUEUE{
		ID: 0x01,
		Parameter: []byte{0x00},
	}
	_, err := xbee.TX(at)
	if err == nil {
		t.Error("Expected error, but got none")
	}

	at = &tx.AT_QUEUE{
		ID: 0x01,
		Command: []byte{'N','I'},
		Parameter: []byte{0x00},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = &tx.AT_QUEUE{
		ID: 0x01,
		Command: []byte{'N','I'},
		Parameter: []byte{0x11},
	}
	xbee.SetApiMode(2)
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_RX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)
	xbee.SetApiMode(2)

	// a valid AT command response
	response := []byte{
		0x7e, 0x00, 0x18, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x20, 0x5a, 0x69, 0x67,
		0x42, 0x65, 0x65, 0x20,
		0x43, 0x6f, 0x6f, 0x72,
		0x64, 0x69, 0x6e, 0x61,
		0x74, 0x6f, 0x72, 0xe5}

	for _, b := range response {
		err := xbee.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}
}

func TestXBee_RX_AT_Escape(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	x := gobee.NewXBee(transmitter, receiver)
	x.SetApiMode(2)

	// a valid AT command response with escaped bytes
	response := []byte{
		0x7e, 0x00, 0x06, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x7D, 0x31, 0xce}

	for _, b := range response {
		err := x.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}
}
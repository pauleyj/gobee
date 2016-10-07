package gobee_test

import (
	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
	"testing"
)

func addressOf(b byte) *byte { return &b }

type Transmitter struct {
	t        *testing.T
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

func (r *Receiver) Receive(f rx.Frame) error {
	switch f.(type) {
	case *rx.AT:
		validateAT(r.t, f.(*rx.AT))

	}

	return nil
}

func validateAT(t *testing.T, f *rx.AT) {
	if f.ID() != 1 {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02x", f.ID)
	}
}

func TestXBee_TX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	var at = tx.NewATBuilder().ID(0x01).Command([2]byte{'N','I'}).Parameter(addressOf(0x00)).Build()
	_, err := xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = tx.NewATBuilder().ID(0x01).Command([2]byte{'N','I'}).Parameter(addressOf(0x11)).Build()
	xbee.SetAPIMode(gobee.EscapeModeActive)
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_AT_QUEUE(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	var at = tx.NewATQueueBuilder().ID(0x01).Command([2]byte{'N','I'}).Parameter(addressOf(0x00)).Build()
	_, err := xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = tx.NewATQueueBuilder().ID(0x01).Command([2]byte{'N','I'}).Parameter(addressOf(0x11)).Build()
	xbee.SetAPIMode(gobee.EscapeModeActive)
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_RX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

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
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func TestXBee_RX_AT_Escape(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewWithEscapeMode(transmitter, receiver, gobee.EscapeModeActive)

	// a valid AT command response with escaped bytes
	response := []byte{
		0x7e, 0x00, 0x06, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x7D, 0x31, 0xce}

	for _, b := range response {
		err := xbee.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func TestXBee_TX_Escape_Length(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewWithEscapeMode(transmitter, receiver, gobee.EscapeModeActive)

	fakeParam := make([]byte, 0)
	for i := 0; i < 0x110D; i++ {
		fakeParam = append(fakeParam, 0)
	}

	zb := tx.NewZBBuilder().ID(0x01).Addr64(0).Addr16(0).BroadcastRadius(0).Options(0).Data(fakeParam).Build()
	_, err := xbee.TX(zb)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_Escape_Checksum(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewWithEscapeMode(transmitter, receiver, gobee.EscapeModeActive)

	at := tx.NewATBuilder().ID(0x01).Command([2]byte{'A','O'}).Parameter(addressOf(0xE8)).Build()
	_, err := xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_Invalid_API_Mode(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	err := xbee.SetAPIMode(gobee.APIEscapeMode(3))
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func TestXBee_Rx_Unknown_Frame_Type(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	unknownFrame := []byte{0x7E, 0x00, 0x10, 0xFF}

	for i, b := range unknownFrame {
		err := xbee.RX(b)
		if i == 3 {
			if err == nil {
				t.Error("Expected error, but got none")
			}
		} else if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func TestXBee_RX_Invalid_Frame_Delimiter(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	err := xbee.RX(0x7D)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func TestXBee_RX_Invalid_Checksum(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	bad_checksum := []byte{
		0x7e, 0x00, 0x0f, 0x97, 0x02, 0x00, 0x13, 0xa2, 0x00,
		0x40, 0x32, 0x03, 0xcf, 0x00, 0x00, 0x41, 0x4f, 0x00,
		0xd0,
	}
	for i, b := range bad_checksum {
		err := xbee.RX(b)
		if i == 18 {
			if err == nil {
				t.Error("Expected error, but got none")
			}
		} else if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

package gobee_test

import (
	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
	"testing"
)

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
	if f.ID != 1 {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02x", f.ID)
	}
}

func TestXBee_TX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	var at = &tx.AT{
		ID:        0x01,
		Parameter: []byte{0x00},
	}
	_, err := xbee.TX(at)
	if err == nil {
		t.Error("Expected error, but got none")
	}

	at = &tx.AT{
		ID:        0x01,
		Command:   []byte{'N', 'I'},
		Parameter: []byte{0x00},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = &tx.AT{
		ID:        0x01,
		Command:   []byte{'N', 'I'},
		Parameter: []byte{0x11},
	}
	xbee.SetAPIMode(2)
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_AT_QUEUE(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	var at = &tx.ATQueue{
		ID:        0x01,
		Parameter: []byte{0x00},
	}
	_, err := xbee.TX(at)
	if err == nil {
		t.Error("Expected error, but got none")
	}

	at = &tx.ATQueue{
		ID:        0x01,
		Command:   []byte{'N', 'I'},
		Parameter: []byte{0x00},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at = &tx.ATQueue{
		ID:        0x01,
		Command:   []byte{'N', 'I'},
		Parameter: []byte{0x11},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_RX_AT(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

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
	xbee := gobee.NewXBee(transmitter, receiver)

	err := xbee.SetAPIMode(2)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// a valid AT command response with escaped bytes
	response := []byte{
		0x7e, 0x00, 0x06, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x7D, 0x31, 0xce}

	for _, b := range response {
		err := xbee.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}
}

func TestXBee_TX_Escape_Length(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	err := xbee.SetAPIMode(2)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	fakeParam := make([]byte, 0)
	for i := 0; i < 0x110D; i++ {
		fakeParam = append(fakeParam, 0)
	}
	at := &tx.AT{
		ID:        0x01,
		Command:   []byte{'A', 'O'},
		Parameter: fakeParam,
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_Escape_Checksum(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	err := xbee.SetAPIMode(2)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	at := &tx.AT{
		ID:        0x01,
		Command:   []byte{'A', 'O'},
		Parameter: []byte{0xE8},
	}
	_, err = xbee.TX(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestXBee_TX_Invalid_API_Mode(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

	err := xbee.SetAPIMode(3)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

const mock_api_id byte = 0xFF

type mock_api_rx_frame struct {
	ID byte
}

func (f *mock_api_rx_frame) RX(b byte) error {
	f.ID = b
	return nil
}

func mockFrameFactoryFunc() rx.Frame {
	return &mock_api_rx_frame{}
}

func TestXBee_Rx_Unknown_Frame_Type(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

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
	xbee := gobee.NewXBee(transmitter, receiver)

	err := xbee.RX(0x7D)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func TestXBee_RX_Invalid_Checksum(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)

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

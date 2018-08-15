package gobee_test

import (
	"fmt"

	"testing"

	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/rx"
	"github.com/pauleyj/gobee/api/tx"
)

type Transmitter struct {
	t        *testing.T
	expected []byte
	i        int
}

func (t *Transmitter) Transmit(p []byte) (int, error) {
	msg := fmt.Sprintf("TX (len = %d)", len(p))
	for _, b := range p {
		msg = fmt.Sprintf("%s 0x%02x ", msg, b)
	}
	t.t.Log(msg)

	if len(t.expected) != len(p) {
		t.t.Fatalf("Expected TX len to be %d, but transmitted %d", len(t.expected), len(p))
	}

	for i, b := range p {
		if b != t.expected[i] {
			t.t.Fatalf("Expected 0x%02x at API frame index %d, but got 0x%02x", t.expected[i], i, b)
		}
	}

	return len(p), nil
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
		t.Errorf("Expected FrameId: 0x01, but got 0x%02x", f.ID())
	}
}

type dummyFrame struct {
	data []byte
}

func (f *dummyFrame) Bytes() ([]byte, error) {
	return f.data, nil
}

type xbeeTXTest struct {
	name     string
	xbeeFunc func(*testing.T) *gobee.XBee
	input    tx.Frame
}

var xbeeTXTests = []xbeeTXTest{
	{
		"TX",
		func(t *testing.T) *gobee.XBee {
			tx := &Transmitter{t: t, expected: []byte{0x7e, 0x00, 0x04, 0x08, 0x01, 0x4e, 0x49, 0x5f}}
			rx := &Receiver{t}
			return gobee.New(tx, rx)
		},
		&dummyFrame{[]byte{0x08, 0x01, 0x4e, 0x49}},
	},
	{
		"TX Escape",
		func(t *testing.T) *gobee.XBee {
			tx := &Transmitter{t: t, expected: []byte{0x7e, 0x00, 0x04, 0x7d, 0x5e, 0x7d, 0x5d, 0x7d, 0x31, 0x7d, 0x33, 0xe0}}
			rx := &Receiver{t}
			return gobee.New(tx, rx, gobee.APIEscapeMode(api.EscapeModeActive))
		},
		&dummyFrame{[]byte{0x7e, 0x7d, 0x11, 0x13}},
	},
}

func TestXBee_TX(t *testing.T) {
	t.Parallel()

	t.Run("XBee TX Suite", func(t *testing.T) {
		for _, tt := range xbeeTXTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				xbee := tt.xbeeFunc(t)
				_, err := xbee.TX(tt.input)
				if err != nil {
					t.Fatalf("Expected no error, but got %v", err)
				}
			})
		}
	})
}

type xbeeRXTest struct {
}

func TestXBee_RX_AT(t *testing.T) {
	transmitter := &Transmitter{t: t}
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
	transmitter := &Transmitter{t: t}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver, gobee.APIEscapeMode(api.EscapeModeActive))

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

//func TestXBee_TX_Escape_Length(t *testing.T) {
//	transmitter := &Transmitter{t: t}
//	receiver := &Receiver{t: t}
//	xbee := gobee.NewWithEscapeMode(transmitter, receiver, api.EscapeModeActive)
//
//	fakeParam := make([]byte, 0)
//	for i := 0; i < 0x110D; i++ {
//		fakeParam = append(fakeParam, 0)
//	}
//
//	zb := zb.NewZB(zb.FrameID(0x01), zb.Addr64(api.BroadcastAddr64), zb.Addr16(api.BroadcastAddr16), zb.BroadcastRadius(0), zb.Data(fakeParam))
//	_, err := xbee.TX(zb)
//	if err != nil {
//		t.Errorf("Expected no error, but got: %v", err)
//	}
//}
//
//func TestXBee_TX_Escape_Checksum(t *testing.T) {
//	transmitter := &Transmitter{t: t}
//	receiver := &Receiver{t: t}
//	xbee := gobee.NewWithEscapeMode(transmitter, receiver, api.EscapeModeActive)
//
//	atFrame := at.NewAT(at.FrameID(0x01), at.Command([2]byte{'A', 'O'}), at.Parameter([]byte{0xE8}))
//	_, err := xbee.TX(atFrame)
//	if err != nil {
//		t.Errorf("Expected no error, but got: %v", err)
//	}
//}
//
//func TestXBee_TX_Invalid_API_Mode(t *testing.T) {
//	transmitter := &Transmitter{t: t}
//	receiver := &Receiver{t: t}
//	xbee := gobee.New(transmitter, receiver)
//
//	err := xbee.SetAPIMode(api.EscapeMode(3))
//	if err == nil {
//		t.Error("Expected error, but got none")
//	}
//}

func TestXBee_Rx_Unknown_Frame_Type(t *testing.T) {
	transmitter := &Transmitter{t: t}
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
	transmitter := &Transmitter{t: t}
	receiver := &Receiver{t: t}
	xbee := gobee.New(transmitter, receiver)

	err := xbee.RX(0x7D)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func TestXBee_RX_Invalid_Checksum(t *testing.T) {
	transmitter := &Transmitter{t: t}
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

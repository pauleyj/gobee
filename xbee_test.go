package gobee_test

import (
	"fmt"
	"testing"
	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/rx"
)

type Transmitter struct{}

func (xm *Transmitter) Write(p []byte) (n int, err error) {
	fmt.Printf("tx <-- %v\n", p)
	return len(p), nil
}

type Receiver struct {
	t *testing.T
}

func (r *Receiver) RxFrameReceiver(f rx.RxFrame) error {
	switch f.(type) {
	case *rx.ATCommandResponse:
		validateAtCommandResponse(r.t, f.(*rx.ATCommandResponse))

	}

	return nil
}

func validateAtCommandResponse(t *testing.T, f *rx.ATCommandResponse) {
	fmt.Printf("%v\n", f)
	if f.FrameId != 1 {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02X", f.FrameId)
	}
}

func TestXBee_RX_AT_Command_Response(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	xbee := gobee.NewXBee(transmitter, receiver)
	xbee.SetApiMode(2)

	// at command response
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

func TestXBee_RX_AT_Command_Response_Escape(t *testing.T) {
	transmitter := &Transmitter{}
	receiver := &Receiver{t: t}
	x := gobee.NewXBee(transmitter, receiver)
	x.SetApiMode(2)

	// at command response
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
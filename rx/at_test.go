package rx

import (
	"bytes"
	"testing"
)

func Test_at_command_response(t *testing.T) {
	// at command response frame data
	response := []byte{
		0x01, 0x4e, 0x49, 0x00,
		0x20, 0x5a, 0x69, 0x67,
		0x42, 0x65, 0x65, 0x20,
		0x43, 0x6f, 0x6f, 0x72,
		0x64, 0x69, 0x6e, 0x61,
		0x74, 0x6f, 0x72, 0xe5}

	rxf := newAT()
	f, ok := rxf.(*AT)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range response {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.FrameId != response[0] {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02X", f.FrameId)
	}

	if !bytes.Equal(f.Command[:1], response[1:2]) {
		t.Errorf("Expected Command: NI, but got %s", string(f.Command[:]))
	}

	if f.Status != response[3] {
		t.Errorf("Expected Status: 0x00, but got 0x%02X", f.Status)
	}

	if !bytes.Equal(f.Data[:], response[4:]) {
		t.Errorf("Expected Data: %v, but got %v", response[4:], f.Data)
	}
}

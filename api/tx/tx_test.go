package tx

import (
	"github.com/pauleyj/gobee/api"
	"testing"
)

func addressOf(b byte) *byte { return &b }

func Test_API_Frame(t *testing.T) {
	at := NewATBuilder().
		ID(0x01).
		Command([2]byte{'N', 'I'}).
		Parameter(nil).
		Build()

	api := &APIFrame{Mode: api.EscapeModeInactive}

	actual, err := api.Bytes(at)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 8 {
		t.Error("Expected length of 8, got %d", len(actual))
	}
}

func Test_API_Frame_WithEscape(t *testing.T) {
	fakeParam := make([]byte, 0)
	for i := 0; i < 0x110D; i++ {
		fakeParam = append(fakeParam, 0)
	}

	zb := NewZBBuilder().
		ID(0x01).
		Addr64(0).
		Addr16(0).
		BroadcastRadius(0).
		Options(0).
		Data(fakeParam).
		Build()

	api := &APIFrame{Mode: api.EscapeModeActive}
	_, err := api.Bytes(zb)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func Test_Valid_AT_No_Param(t *testing.T) {
	at := NewATBuilder().
		ID(0x01).
		Command([2]byte{'N', 'I'}).
		Parameter(nil).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 4 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atAPIID, 1, 'N', 'I'}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_With_Param(t *testing.T) {
	at := NewATBuilder().
		ID(0x01).
		Command([2]byte{'N', 'I'}).
		Parameter(addressOf(1)).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 5 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atAPIID, 1, 'N', 'I', 1}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_REMOTE_No_Param(t *testing.T) {
	at := NewATRemoteBuilder().
		ID(0x01).
		Addr64(0x000000000000FFFF).
		Addr16(0xFFFE).
		Options(0x00).
		Command([2]byte{'A', 'O'}).
		Parameter(nil).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 15 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'A', 'O'}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_REMOTE_With_Param(t *testing.T) {
	at := NewATRemoteBuilder().
		ID(0x01).
		Addr64(0x000000000000FFFF).
		Addr16(0xFFFE).
		Options(0x00).
		Command([2]byte{'A', 'O'}).
		Parameter(addressOf(0x01)).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 16 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'A', 'O', 0x01}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_QUEUE_No_Param(t *testing.T) {
	at := NewATQueueBuilder().
		ID(0x01).
		Command([2]byte{'N', 'I'}).
		Parameter(nil).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 4 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atQueueAPIID, 1, 'N', 'I'}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_QUEUE_With_Param(t *testing.T) {
	at := NewATQueueBuilder().
		ID(0x01).
		Command([2]byte{'N', 'I'}).
		Parameter(addressOf(0x00)).
		Build()

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 5 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{atQueueAPIID, 1, 'N', 'I', 0}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_ZB(t *testing.T) {
	zb := NewZBBuilder().
		ID(0xFF).
		Addr64(0x0001020304050607).
		Addr16(0x0001).
		BroadcastRadius(0xFF).
		Options(0xEE).
		Data([]byte{0x00, 0x01}).
		Build()

	actual, err := zb.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 16 {
		t.Errorf("Expected ZB frame to be 16 bytes in length, got: %d", len(actual))
	}

	expected := []byte{zbAPIID, 0xFF, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
		0x06, 0x07, 0x00, 0x01, 0xFF, 0xEE, 0x00, 0x01}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_ZB_EXPLICIT(t *testing.T) {
	zb := NewZBExplicitBuilder().
		ID(0xFF).
		Addr64(0x0001020304050607).
		Addr16(0x0001).
		SrcEP(0x01).
		DstEP(0x02).
		ClusterID(0x1234).
		ProfileID(0x5678).
		BroadcastRadius(0xFF).
		Options(0xEE).
		Data([]byte{0x00, 0x01}).
		Build()

	actual, err := zb.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 22 {
		t.Errorf("Expected ZB frame to be 16 bytes in length, got: %d", len(actual))
	}

	expected := []byte{zbExplicitAPIID, 0xFF, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
		0x06, 0x07, 0x00, 0x01,
		0x01, 0x02, 0x12, 0x34, 0x56, 0x78,
		0xFF, 0xEE, 0x00, 0x01}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
		}
	}
}

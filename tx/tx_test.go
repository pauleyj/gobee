package tx

import "testing"

func Test_Invalid_AT(t *testing.T) {
	var at = &AT{
		ID: 0x01,
		Parameter: []byte{0x00},
	}

	_, err := at.Bytes()
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func Test_Valid_AT_No_Param(t *testing.T) {
	at := &AT{
		ID: 0x01,
		Command: []byte{'N', 'I'},
	}

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 4 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{api_id_at, 1, 'N', 'I'}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x02%x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_With_Param(t *testing.T) {
	at := &AT{
		ID: 0x01,
		Command: []byte{'N', 'I'},
		Parameter: []byte{0x00},
	}

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 5 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{api_id_at, 1, 'N', 'I', 0}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x02%x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Invalid_AT_QUEUE(t *testing.T) {
	var at = &AT_QUEUE{
		ID: 0x01,
		Parameter: []byte{0x00},
	}

	_, err := at.Bytes()
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

func Test_Valid_AT_QUEUE_No_Param(t *testing.T) {
	at := &AT_QUEUE{
		ID: 0x01,
		Command: []byte{'N', 'I'},
	}

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 4 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{api_id_at_queue, 1, 'N', 'I'}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x02%x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_Valid_AT_QUEUE_With_Param(t *testing.T) {
	at := &AT_QUEUE{
		ID: 0x01,
		Command: []byte{'N', 'I'},
		Parameter: []byte{0x00},
	}

	actual, err := at.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 5 {
		t.Errorf("Expected AT frame to be 5 bytes in length, got: %d", len(actual))
	}

	expected := []byte{api_id_at_queue, 1, 'N', 'I', 0}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x02%x, but got 0x%02x", b, actual[i])
		}
	}
}

func Test_ZB(t *testing.T) {
	zb := &ZB{
		ID: 0xFF,
		Addr64: 0x0001020304050607,
		Addr16: 0x0001,
		BroadcastRadius: 0xFF,
		Options: 0xEE,
		Data: []byte{0x00, 0x01},
	}

	actual, err := zb.Bytes()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if len(actual) != 16 {
		t.Errorf("Expected ZB frame to be 15 bytes in length, got: %d", len(actual))
	}

	expected := []byte{api_id_tx_zb, 0xFF, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
		0x06, 0x07, 0x00, 0x01, 0xFF, 0xEE, 0x00, 0x01}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected 0x02%x, but got 0x%02x", b, actual[i])
		}
	}
}
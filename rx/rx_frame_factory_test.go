package rx

import "testing"

const unknown_api_id byte = 0x00

func TestNewRxFrameForApiId(t *testing.T) {
	rxf, err := NewRxFrameForApiId(XBEE_API_ID_AT_COMMAND_RESPONSE)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_, ok := rxf.(*ATCommandResponse)
	if !ok {
		t.Error("Failed type assertion ATCommandResponse")
	}

	_, err = NewRxFrameForApiId(unknown_api_id)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", errUnknownFrameApiId)
	}
	if err != errUnknownFrameApiId {
		t.Errorf("Expected error: %v, but got: %v", errUnknownFrameApiId, err)
	}
}

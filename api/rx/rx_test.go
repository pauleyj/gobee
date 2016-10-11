package rx

import (
	"bytes"
	"github.com/pauleyj/gobee/api"
	"testing"
)

func Test_API_Frame(t *testing.T) {
	// a valid AT command response
	response := []byte{
		0x7e, 0x00, 0x18, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x20, 0x5a, 0x69, 0x67,
		0x42, 0x65, 0x65, 0x20,
		0x43, 0x6f, 0x6f, 0x72,
		0x64, 0x69, 0x6e, 0x61,
		0x74, 0x6f, 0x72, 0xe5}

	api := &APIFrame{
		Mode: api.EscapeModeInactive,
	}

	for _, c := range response {
		f, err := api.RX(c)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if f != nil {
			_, ok := f.(*AT)
			if !ok {
				t.Error("Failed type assertion")
			}
		}
	}
}

func Test_API_Frame_With_Escape(t *testing.T) {
	// a valid AT command response with escaped bytes
	response := []byte{
		0x7e, 0x00, 0x06, 0x88,
		0x01, 0x4e, 0x49, 0x00,
		0x7D, 0x31, 0xce}

	api := &APIFrame{
		Mode: api.EscapeModeActive,
	}

	for _, c := range response {
		f, err := api.RX(c)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if f != nil {
			_, ok := f.(*AT)
			if !ok {
				t.Error("Failed type assertion")
			}
		}
	}
}

func Test_API_Frame_Invalid_Checksum(t *testing.T) {
	bad_checksum := []byte{
		0x7e, 0x00, 0x0f, 0x97, 0x02, 0x00, 0x13, 0xa2, 0x00,
		0x40, 0x32, 0x03, 0xcf, 0x00, 0x00, 0x41, 0x4f, 0x00,
		0xd0,
	}
	api := &APIFrame{
		Mode: api.EscapeModeActive,
	}

	for i, c := range bad_checksum {
		_, err := api.RX(c)
		if i == 18 {

			if err == nil {
				t.Error("Expected error, but got none")
			}
		} else if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}
}

func Test_AT(t *testing.T) {
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

	if f.ID() != response[0] {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02X", f.ID())
	}

	if !bytes.Equal(f.Command(), []byte{'N', 'I'}) {
		t.Errorf("Expected Command: NI, but got %s", f.Command())
	}

	if f.Status() != response[3] {
		t.Errorf("Expected Status: 0x00, but got 0x%02X", f.Status())
	}

	if !bytes.Equal(f.Data(), response[4:]) {
		t.Errorf("Expected Data: %v, but got %v", response[4:], f.Data())
	}
}

func Test_ZB(t *testing.T) {
	// zb frame data
	actual := []byte{0x00, 0x13, 0xa2, 0x00, 0x40, 0x32, 0x03, 0xab,
		0x5f, 0xd6,
		0x01,
		0x66, 0x6f, 0x6f}

	rxf := newZB()
	f, ok := rxf.(*ZB)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.Addr64() != 0x0013A200403203AB {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64)
	}

	if f.Addr16() != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16)
	}

	if f.Options() != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options)
	}

	if !bytes.Equal(f.Data(), []byte{'f', 'o', 'o'}) {
		t.Errorf("Expected Data: %v, but got %v", []byte{'f', 'o', 'o'}, f.Data)
	}
}

func Test_ZB_Explicit(t *testing.T) {
	// zb explicit frame data
	actual := []byte{
		0x00, 0x13, 0xa2, 0x00, 0x40, 0x32, 0x03, 0xab,
		0x5f, 0xd6,
		0xcd,
		0x01,
		0x00, 0x54,
		0xc1, 0x05,
		0x01,
		0x66, 0x6f, 0x6f}

	rxf := newZBExplicit()
	f, ok := rxf.(*ZBExplicit)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.Addr64() != 0x0013A200403203AB {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64)
	}

	if f.Addr16() != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16)
	}

	if f.SrcEP() != 0xCD {
		t.Errorf("Expected SrcEP to be 0x%02X, but got 0x%02X", 0xCD, f.SrcEP)
	}

	if f.DstEP() != 0x01 {
		t.Errorf("Expected DstEP to be 0x%02X, but got 0x%02X", 0x01, f.DstEP)
	}

	if f.ClusterID() != 0x0054 {
		t.Errorf("Expected ClusterID to be 0x%04X, but got 0x%04X", 0x54C1, f.ClusterID)
	}

	if f.ProfileID() != 0xC105 {
		t.Errorf("Expected ProfileID to be 0x%04X, but got 0x%04X", 0x0501, f.ProfileID)
	}

	if f.Options() != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options)
	}

	if !bytes.Equal(f.Data(), []byte{'f', 'o', 'o'}) {
		t.Errorf("Expected Data: %v, but got %v", []byte{'f', 'o', 'o'}, f.Data)
	}
}

func Test_TX_STATUS(t *testing.T) {
	actual := []byte{
		0x01,
		0xff, 0xfe,
		0x00,
		0x00,
		0x00,
	}

	rxf := newTXStatus()
	f, ok := rxf.(*TXStatus)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.ID() != 0x01 {
		t.Errorf("Expected ID = 0x%02X, but got 0x%02X", 0x01, f.ID())
	}

	if f.Addr16() != 0xFFFE {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0xFFFE, f.Addr16())
	}

	if f.Retries() != 0x00 {
		t.Errorf("Expected Retries = 0x%02X, but got 0x%02X", 0x01, f.Retries())
	}

	if f.Delivery() != 0x00 {
		t.Errorf("Expected Delivery = 0x%02X, but got 0x%02X", 0x01, f.Delivery())
	}

	if f.Discovery() != 0x00 {
		t.Errorf("Expected Discovery = 0x%02X, but got 0x%02X", 0x01, f.Discovery())
	}
}

func Test_AT_REMOTE(t *testing.T) {
	actual := []byte{
		0x01,
		0x00, 0x13, 0xa2, 0x00, 0x40, 0x32, 0x03, 0xcf,
		0x00, 0x00,
		0x41, 0x4f,
		0x00,
		0x02,
	}

	rxf := newATRemote()
	f, ok := rxf.(*ATRemote)
	if !ok {
		t.Error("Failed type assertion AT_REMOTE")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.ID() != 0x01 {
		t.Errorf("Expected ID = 0x%02X, but got 0x%02X", 0x01, f.ID())
	}

	if f.Addr64() != 0x0013a200403203cf {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013a200403203cf, f.Addr64())
	}

	if f.Addr16() != 0x0000 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x0000, f.Addr16())
	}

	if !bytes.Equal(f.Command(), []byte{'A', 'O'}) {
		t.Errorf("Expected command to be AO, but got %s", string(f.Command()))
	}

	if f.Status() != 0x00 {
		t.Errorf("Expected Status = 0x%02X, but got 0x%02X", 0x00, f.Status())
	}

	if len(f.Data()) != 0x01 {
		t.Errorf("Expected Data length to be 0x%02X, but is 0x%02X", 0x01, len(f.Data()))
	}

	if f.Data()[0] != 0x02 {
		t.Errorf("Expected Data to be 0x%02X, but got 0x%02X", 0x02, f.Data()[0])
	}
}

const unknownAPIID byte = 0x00

func TestNewRxFrameForApiId(t *testing.T) {
	rxf, err := NewFrameForAPIID(atAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok := rxf.(*AT)
	if !ok {
		t.Error("Failed type assertion AT")
	}

	rxf, err = NewFrameForAPIID(zbAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*ZB)
	if !ok {
		t.Error("Failed type assertion ZB")
	}

	rxf, err = NewFrameForAPIID(txStatusAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*TXStatus)
	if !ok {
		t.Error("Failed type assertion TX_STATUS")
	}

	rxf, err = NewFrameForAPIID(zbExplicitAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*ZBExplicit)
	if !ok {
		t.Error("Failed type assertion ZB_EXPLICIT")
	}

	rxf, err = NewFrameForAPIID(atRemoteAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*ATRemote)
	if !ok {
		t.Error("Failed type assertion AT_REMOTE")
	}

	_, err = NewFrameForAPIID(unknownAPIID)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", errUnknownFrameAPIID)
	}
	if err != errUnknownFrameAPIID {
		t.Errorf("Expected error: %v, but got: %v", errUnknownFrameAPIID, err)
	}
}

const mockAPIID byte = 0xFF

type mockFrame struct {
	ID byte
}

func (f *mockFrame) RX(b byte) error {
	f.ID = b
	return nil
}

func mockFrameFactoryFunc() Frame {
	return &mockFrame{}
}

func TestAddNewAPIFrameFactory(t *testing.T) {
	err := AddFactoryForAPIID(mockAPIID, mockFrameFactoryFunc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rxf, err := NewFrameForAPIID(mockAPIID)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	_, ok := rxf.(*mockFrame)
	if !ok {
		t.Error("Failed type assertion mock_api_rx_frame")
	}

}

func TestAddExistingAPIFrameFactory(t *testing.T) {
	err := AddFactoryForAPIID(atAPIID, mockFrameFactoryFunc)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

package rx

import (
	"bytes"
	"testing"
)

const unknown_api_id byte = 0x00

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

	if f.ID != response[0] {
		t.Errorf("Expected FrameId: 0x01, but got 0x%02X", f.ID)
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

	if f.Addr64 != 0x0013A200403203AB {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64)
	}

	if f.Addr16 != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16)
	}

	if f.Options != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options)
	}

	if !bytes.Equal(f.Data[:], []byte{'f', 'o', 'o'}) {
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

	rxf := newZB_EXPLICIT()
	f, ok := rxf.(*ZB_EXPLICIT)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.Addr64 != 0x0013A200403203AB {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64)
	}

	if f.Addr16 != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16)
	}

	if f.SrcEP != 0xCD {
		t.Errorf("Expected SrcEP to be 0x%02X, but got 0x%02X", 0xCD, f.SrcEP)
	}

	if f.DstEP != 0x01 {
		t.Errorf("Expected DstEP to be 0x%02X, but got 0x%02X", 0x01, f.DstEP)
	}

	if f.ClusterID != 0x0054 {
		t.Errorf("Expected ClusterID to be 0x%04X, but got 0x%04X", 0x54C1, f.ClusterID)
	}

	if f.ProfileID != 0xC105 {
		t.Errorf("Expected ProfileID to be 0x%04X, but got 0x%04X", 0x0501, f.ProfileID)
	}

	if f.Options != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options)
	}

	if !bytes.Equal(f.Data[:], []byte{'f', 'o', 'o'}) {
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

	rxf := newTX_STATUS()
	f, ok := rxf.(*TX_STATUS)
	if !ok {
		t.Error("Failed type assertion")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.ID != 0x01 {
		t.Errorf("Expected ID = 0x%02X, but got 0x%02X", 0x01, f.ID)
	}

	if f.Addr16 != 0xFFFE {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0xFFFE, f.Addr16)
	}

	if f.Retries != 0x00 {
		t.Errorf("Expected Retries = 0x%02X, but got 0x%02X", 0x01, f.ID)
	}

	if f.Delivery != 0x00 {
		t.Errorf("Expected Delivery = 0x%02X, but got 0x%02X", 0x01, f.ID)
	}

	if f.Discovery != 0x00 {
		t.Errorf("Expected Discovery = 0x%02X, but got 0x%02X", 0x01, f.ID)
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

	rxf := newAT_REMOTE()
	f, ok := rxf.(*AT_REMOTE)
	if !ok {
		t.Error("Failed type assertion AT_REMOTE")
	}

	for _, b := range actual {
		err := f.RX(b)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
	}

	if f.ID != 0x01 {
		t.Errorf("Expected ID = 0x%02X, but got 0x%02X", 0x01, f.ID)
	}

	if f.Addr64 != 0x0013a200403203cf {
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013a200403203cf, f.Addr64)
	}

	if f.Addr16 != 0x0000 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x0000, f.Addr16)
	}

	if !bytes.Equal(f.Command[:], []byte{'A', 'O'}) {
		t.Errorf("Expected command to be AO, but got %s", string(f.Command[:]))
	}

	if f.Status != 0x00 {
		t.Errorf("Expected Status = 0x%02X, but got 0x%02X", 0x00, f.ID)
	}

	if len(f.Data) != 0x01 {
		t.Errorf("Expected Data length to be 0x%02X, but is 0x%02X", 0x01, len(f.Data))
	}

	if f.Data[0] != 0x02 {
		t.Errorf("Expected Data to be 0x%02X, but got 0x%02X", 0x02, f.Data[0])
	}
}

func TestNewRxFrameForApiId(t *testing.T) {
	rxf, err := NewRxFrameForApiId(XBEE_API_ID_RX_AT)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok := rxf.(*AT)
	if !ok {
		t.Error("Failed type assertion AT")
	}

	rxf, err = NewRxFrameForApiId(XBEE_API_ID_RX_ZB)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*ZB)
	if !ok {
		t.Error("Failed type assertion ZB")
	}

	rxf, err = NewRxFrameForApiId(XBEE_API_TX_STATUS)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*TX_STATUS)
	if !ok {
		t.Error("Failed type assertion TX_STATUS")
	}

	rxf, err = NewRxFrameForApiId(XBEE_API_ID_RX_ZB_EXPLICIT)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*ZB_EXPLICIT)
	if !ok {
		t.Error("Failed type assertion ZB_EXPLICIT")
	}

	rxf, err = NewRxFrameForApiId(XBEE_API_ID_RX_AT_REMOTE)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	_, ok = rxf.(*AT_REMOTE)
	if !ok {
		t.Error("Failed type assertion AT_REMOTE")
	}

	_, err = NewRxFrameForApiId(unknown_api_id)
	if err == nil {
		t.Errorf("Expected error: %v, but got none", errUnknownFrameApiId)
	}
	if err != errUnknownFrameApiId {
		t.Errorf("Expected error: %v, but got: %v", errUnknownFrameApiId, err)
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

func mockFrameFactoryFunc() RxFrame {
	return &mock_api_rx_frame{}
}

func TestAddNewAPIFrameFactory(t *testing.T) {
	err := AddApiFactoryForId(mock_api_id, mockFrameFactoryFunc)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	rxf, err := NewRxFrameForApiId(mock_api_id)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	_, ok := rxf.(*mock_api_rx_frame)
	if !ok {
		t.Error("Failed type assertion mock_api_rx_frame")
	}

}

func TestAddExistingAPIFrameFactory(t *testing.T) {
	err := AddApiFactoryForId(XBEE_API_ID_RX_AT, mockFrameFactoryFunc)
	if err == nil {
		t.Error("Expected error, but got none")
	}
}

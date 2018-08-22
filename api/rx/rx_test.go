package rx

import (
	"bytes"
	"github.com/pauleyj/gobee/api"
	"testing"
	"reflect"
	"errors"
)

type rxTest struct {
	name     string
	input    []byte
	f        *APIFrame
	expected Frame
	err      error
}

var badFrameError = errors.New("bad frame")
const badFrameAPIID = byte(0xfe)
type badFrame struct{}

func (f *badFrame) RX(c byte) error {
	return badFrameError
}

var rxTests = []rxTest{
	{"Checksum Error",
		[]byte{
			0x7e, 0x00, 0x0f, 0x97, 0x02, 0x00, 0x13, 0xa2, 0x00,
			0x40, 0x32, 0x03, 0xcf, 0x00, 0x00, 0x41, 0x4f, 0x00,
			0xd0,
		},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		nil,
		api.ErrChecksumValidation,
	},
	{"Unknown RX Frame ID",
		[]byte{0x7e, 0x00, 0x18, 0xff},
		New(),
		nil,
		errUnknownFrameAPIID,
	},
	{"Start Frame Error",
		[]byte{0x0f},
		New(),
		nil,
		api.ErrFrameDelimiter,
	},
	{"Bad Frame",
		[]byte{0x7e, 0x00, 0x18, badFrameAPIID, 0x00},
		New(),
		nil,
		badFrameError,
	},
	{"RX AT Frame",
		[]byte{
			0x7e, 0x00, 0x18, 0x88,
			0x01, 0x4e, 0x49, 0x00,
			0x20, 0x5a, 0x69, 0x67,
			0x42, 0x65, 0x65, 0x20,
			0x43, 0x6f, 0x6f, 0x72,
			0x64, 0x69, 0x6e, 0x61,
			0x74, 0x6f, 0x72, 0xe5},
		New(),
		&AT{[]byte{
			0x01, 0x4e, 0x49, 0x00,
			0x20, 0x5a, 0x69, 0x67,
			0x42, 0x65, 0x65, 0x20,
			0x43, 0x6f, 0x6f, 0x72,
			0x64, 0x69, 0x6e, 0x61,
			0x74, 0x6f, 0x72}},
		nil,
	},
	{"RX AT Frame with Escape",
		[]byte{
			0x7e, 0x00, 0x06, 0x88,
			0x01, 0x4e, 0x49, 0x00,
			0x7D, 0x31, 0xce},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		&AT{[]byte{0x01, 0x4e, 0x49, 0x00, 0x11}},
		nil,
	},
	{"RX API Frame nil Options",
		[]byte{
			0x7e, 0x00, 0x18, 0x88,
			0x01, 0x4e, 0x49, 0x00,
			0x20, 0x5a, 0x69, 0x67,
			0x42, 0x65, 0x65, 0x20,
			0x43, 0x6f, 0x6f, 0x72,
			0x64, 0x69, 0x6e, 0x61,
			0x74, 0x6f, 0x72, 0xe5},
		New(nil),
		&AT{[]byte{
			0x01, 0x4e, 0x49, 0x00,
			0x20, 0x5a, 0x69, 0x67,
			0x42, 0x65, 0x65, 0x20,
			0x43, 0x6f, 0x6f, 0x72,
			0x64, 0x69, 0x6e, 0x61,
			0x74, 0x6f, 0x72}},
		nil,
	},
}

func TestRXAPIFrame(t *testing.T) {
	// register phony API frame factory
	AddFactoryForAPIID(badFrameAPIID, func() Frame {
		return &badFrame{}
	})

	t.Parallel()

	t.Run("RX APIFrame Suite", func(t *testing.T) {
		for _, tt := range rxTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				var actual Frame
				var err error
				for _, c := range tt.input {
					actual, err = tt.f.RX(c)
					if tt.err != nil && err != nil {
						if err != tt.err {
							t.Fatalf("Expected error=%+v, but got %+v", tt.err, err)
						}
					} else {
						if err != nil {
							t.Fatalf("expected no error, but got %v", err)
						}
					}
				}

				expectedType := reflect.TypeOf(tt.expected)
				actualType := reflect.TypeOf(actual)
				if actualType != expectedType {
					t.Fatalf("expected type %v, but got %v", expectedType, actualType)
				}

				if f, ok := actual.(IDGetter); ok {
					e := tt.expected.(IDGetter)
					if f.ID() != e.ID() {
						t.Fatalf("Expected frame ID=0x%02x, but got 0x%02x", e.ID(), f.ID())
					}
				}

				if f, ok := actual.(Addr64Getter); ok {
					e := tt.expected.(Addr64Getter)
					if f.Addr64() != e.Addr64() {
						t.Fatalf("Expected frame Addr64=0x%08x, but got 0x%08x", e.Addr64(), f.Addr64())
					}
				}

				if f, ok := actual.(Addr16Getter); ok {
					e := tt.expected.(Addr16Getter)
					if f.Addr16() != e.Addr16() {
						t.Fatalf("Expected frame Addr16=0x%04x, but got 0x%04x", e.Addr16(), f.Addr16())
					}
				}

				if f, ok := actual.(SrcEPGetter); ok {
					e := tt.expected.(SrcEPGetter)
					if f.SrcEP() != e.SrcEP() {
						t.Fatalf("Expected frame SrcEP=0x%02x, but got 0x%02x", e.SrcEP(), f.SrcEP())
					}
				}

				if f, ok := actual.(DstEPGetter); ok {
					e := tt.expected.(DstEPGetter)
					if f.DstEP() != e.DstEP() {
						t.Fatalf("Expected frame DstEP=0x%02x, but got 0x%02x", e.DstEP(), f.DstEP())
					}
				}

				if f, ok := actual.(ClusterIDGetter); ok {
					e := tt.expected.(ClusterIDGetter)
					if f.ClusterID() != e.ClusterID() {
						t.Fatalf("Expected frame ClusterID=0x%04x, but got 0x%04x", e.ClusterID(), f.ClusterID())
					}
				}

				if f, ok := actual.(ProfileIDGetter); ok {
					e := tt.expected.(ProfileIDGetter)
					if f.ProfileID() != e.ProfileID() {
						t.Fatalf("Expected frame ProfileID=0x%04x, but got 0x%04x", e.ProfileID(), f.ProfileID())
					}
				}

				if f, ok := actual.(CommandGetter); ok {
					e := tt.expected.(CommandGetter)
					ac := f.Command()
					ec := e.Command()

					if len(ac) != len(ec) {
						t.Fatalf("Expected len(command)=%d, but got %d", len(ec), len(ac))
					}

					for i, c := range ec {
						if ac[i] != c {
							t.Fatalf("Expected command[%d]=0x%02x, but got 0x%02x", i, c, ac[i])
						}
					}
				}

				if f, ok := actual.(StatusGetter); ok {
					e := tt.expected.(StatusGetter)
					if f.Status() != e.Status() {
						t.Fatalf("Expected frame Status=0x%02x, but got 0x%02x", e.Status(), f.Status())
					}
				}

				if f, ok := actual.(OptionsGetter); ok {
					e := tt.expected.(OptionsGetter)
					if f.Options() != e.Options() {
						t.Fatalf("Expected frame Options=0x%02x, but got 0x%02x", e.Options(), f.Options())
					}
				}

				if f, ok := actual.(DataGetter); ok {
					e := tt.expected.(DataGetter)
					ad := f.Data()
					ed := e.Data()

					if len(ad) != len(ed) {
						t.Fatalf("Expected len(data)=%d, but got %d", len(ed), len(ad))
					}

					for i, c := range ed {
						if ad[i] != c {
							t.Fatalf("Expected data[%d]=0x%02x, but got 0x%02x", i, c, ad[i])
						}
					}
				}

				if f, ok := actual.(RetriesGetter); ok {
					e := tt.expected.(RetriesGetter)
					if f.Retries() != e.Retries() {
						t.Fatalf("Expected frame Retries=0x%02x, but got 0x%02x", e.Retries(), f.Retries())
					}
				}

				if f, ok := actual.(DeliveryGetter); ok {
					e := tt.expected.(DeliveryGetter)
					if f.Delivery() != e.Delivery() {
						t.Fatalf("Expected frame Delivery=0x%02x, but got 0x%02x", e.Delivery(), f.Delivery())
					}
				}

				if f, ok := actual.(DiscoveryGetter); ok {
					e := tt.expected.(DiscoveryGetter)
					if f.Discovery() != e.Discovery() {
						t.Fatalf("Expected frame Discovery=0x%02x, but got 0x%02x", e.Discovery(), f.Discovery())
					}
				}
			})
		}
	})
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
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64())
	}

	if f.Addr16() != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16())
	}

	if f.Options() != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options())
	}

	if !bytes.Equal(f.Data(), []byte{'f', 'o', 'o'}) {
		t.Errorf("Expected Data: %v, but got %v", []byte{'f', 'o', 'o'}, f.Data())
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
		t.Errorf("Expected Addr64 to be 0x%016X, but got 0x%016X", 0x0013A200403203AB, f.Addr64())
	}

	if f.Addr16() != 0x5FD6 {
		t.Errorf("Expected Addr16 to be 0x%04X, but got 0x%04X", 0x5FD6, f.Addr16())
	}

	if f.SrcEP() != 0xCD {
		t.Errorf("Expected SrcEP to be 0x%02X, but got 0x%02X", 0xCD, f.SrcEP())
	}

	if f.DstEP() != 0x01 {
		t.Errorf("Expected DstEP to be 0x%02X, but got 0x%02X", 0x01, f.DstEP())
	}

	if f.ClusterID() != 0x0054 {
		t.Errorf("Expected ClusterID to be 0x%04X, but got 0x%04X", 0x54C1, f.ClusterID())
	}

	if f.ProfileID() != 0xC105 {
		t.Errorf("Expected ProfileID to be 0x%04X, but got 0x%04X", 0x0501, f.ProfileID())
	}

	if f.Options() != 0x01 {
		t.Errorf("Expected Options to be 0x%02X, but got 0x%02X", 0x01, f.Options())
	}

	if !bytes.Equal(f.Data(), []byte{'f', 'o', 'o'}) {
		t.Errorf("Expected Data: %v, but got %v", []byte{'f', 'o', 'o'}, f.Data())
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

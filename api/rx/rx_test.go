package rx

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pauleyj/gobee/api"
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
	//{"Bad Frame",
	//	[]byte{0x7e, 0x00, 0x18, badFrameAPIID, 0x00},
	//	New(),
	//	nil,
	//	badFrameError,
	//},
	{"RX AT Frame With Data",
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
	{
		name:     "RX AT Frame Without Data",
		input:    []byte{0x7e, 0x00, 0x05, 0x88, 0x01, 0x42, 0x44, 0x00, 0xf0},
		f:        New(),
		expected: &AT{[]byte{0x01, 0x42, 0x44, 0x00}},
		err:      nil,
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
	{
		name:     "RX Modem Status",
		input:    []byte{0x7e, 0x00, 0x02, 0x8A, 0x06, 0x6F},
		f:        New(),
		expected: &ModemStatus{0x06},
		err:      nil,
	},
	{
		name: "RX Remote Command Response Without Data",
		input: []byte{
			0x7e, 0x00, 0x0f, 0x97,
			0x55, 0x00, 0x13, 0xa2,
			0x00, 0x40, 0x52, 0x2b,
			0xaa, 0x7d, 0x84, 0x53,
			0x4c, 0x00, 0x57},
		f: New(),
		expected: &ATRemote{
			[]byte{
				0x55, 0x00, 0x13, 0xa2,
				0x00, 0x40, 0x52, 0x2b,
				0xaa, 0x7d, 0x84, 0x53,
				0x4c, 0x00},
		},
		err: nil,
	},
	{
		name: "RX Remote Command Response With Data",
		input: []byte{
			0x7e, 0x00, 0x13, 0x97,
			0x55, 0x00, 0x13, 0xa2,
			0x00, 0x40, 0x52, 0x2b,
			0xaa, 0x7d, 0x84, 0x53,
			0x4c, 0x00, 0x40, 0x52,
			0x2b, 0xaa, 0xf0},
		f: New(),
		expected: &ATRemote{
			[]byte{
				0x55, 0x00, 0x13, 0xa2,
				0x00, 0x40, 0x52, 0x2b,
				0xaa, 0x7d, 0x84, 0x53,
				0x4c, 0x00, 0x40, 0x52,
				0x2b, 0xaa},
		},
		err: nil,
	},
	{
		name: "RX ZB",
		input: []byte{
			0x7E, 0x00, 0x12, 0x90,
			0x00, 0x13, 0xA2, 0x00,
			0x40, 0x52, 0x2B, 0xAA,
			0x7D, 0x84, 0x01, 0x52,
			0x78, 0x44, 0x61, 0x74,
			0x61, 0x0D},
		f: New(),
		expected: &ZB{
			[]byte{
				0x00, 0x13, 0xA2, 0x00,
				0x40, 0x52, 0x2B, 0xAA,
				0x7D, 0x84, 0x01, 0x52,
				0x78, 0x44, 0x61, 0x74,
				0x61},
		},
		err: nil,
	},
	{
		name: "RX ZB Without Data",
		input: []byte{
			0x7E, 0x00, 0x0c, 0x90,
			0x00, 0x13, 0xA2, 0x00,
			0x40, 0x52, 0x2B, 0xAA,
			0x7D, 0x84, 0x01, 0x51},
		f: New(),
		expected: &ZB{
			[]byte{
				0x00, 0x13, 0xA2, 0x00,
				0x40, 0x52, 0x2B, 0xAA,
				0x7D, 0x84, 0x01},
		},
		err: nil,
	},
	{
		name: "RX ZB Explicit",
		input: []byte{
			0x7E, 0x00, 0x18, 0x91,
			0x00, 0x13, 0xA2, 0x00,
			0x40, 0x52, 0x2B, 0xAA,
			0x7D, 0x84, 0xE0, 0xE0,
			0x22, 0x11, 0xC1, 0x05,
			0x02, 0x52, 0x78, 0x44,
			0x61, 0x74, 0x61, 0x52},
		f: New(),
		expected: &ZBExplicit{
			[]byte{
				0x00, 0x13, 0xA2, 0x00,
				0x40, 0x52, 0x2B, 0xAA,
				0x7D, 0x84, 0xE0, 0xE0,
				0x22, 0x11, 0xC1, 0x05,
				0x02, 0x52, 0x78, 0x44,
				0x61, 0x74, 0x61},
		},
		err: nil,
	},
	{
		name: "RX ZB Explicit No Data",
		input: []byte{
			0x7E, 0x00, 0x12, 0x91,
			0x00, 0x13, 0xA2, 0x00,
			0x40, 0x52, 0x2B, 0xAA,
			0x7D, 0x84, 0xE0, 0xE0,
			0x22, 0x11, 0xC1, 0x05,
			0x02, 0x96},
		f: New(),
		expected: &ZBExplicit{
			[]byte{
				0x00, 0x13, 0xA2, 0x00,
				0x40, 0x52, 0x2B, 0xAA,
				0x7D, 0x84, 0xE0, 0xE0,
				0x22, 0x11, 0xC1, 0x05,
				0x02},
		},
		err: nil,
	},
	{
		name: "RX TX Status",
		input: []byte{
			0x7E, 0x00, 0x07, 0x8B,
			0x01, 0x7D, 0x84, 0x00,
			0x00, 0x01, 0x71},
		f: New(),
		expected: &TXStatus{
			[]byte{
				0x01, 0x7D, 0x84, 0x00,
				0x00, 0x01},
		},
		err: nil,
	},
	{
		name: "RX IO Data Sample",
		input: []byte{
			0x7E, 0x00, 0x14, 0x92,
			0x00, 0x13, 0xA2, 0x00,
			0x40, 0x52, 0x2B, 0xAA,
			0x7D, 0x84, 0x01, 0x01,
			0x00, 0x1C, 0x02, 0x00,
			0x14, 0x02, 0x25, 0xF5},
		f: New(),
		expected: &IOSample{
			[]byte{
				0x00, 0x13, 0xA2, 0x00,
				0x40, 0x52, 0x2B, 0xAA,
				0x7D, 0x84, 0x01, 0x01,
				0x00, 0x1C, 0x02, 0x00,
				0x14, 0x02, 0x25},
		},
		err: nil,
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

				if f, ok := actual.(SampleCountGetter); ok {
					e := tt.expected.(SampleCountGetter)
					if f.SampleCount() != e.SampleCount() {
						t.Fatalf("Expected sample count=%#0.2x, but got %#0.2x", e.SampleCount(), f.SampleCount())
					}
				}

				if f, ok := actual.(DigitalSampleMaskGetter); ok {
					e := tt.expected.(DigitalSampleMaskGetter)
					if f.DigitalSampleMask() != e.DigitalSampleMask() {
						t.Fatalf("Expected digital sample mask=%#0.4x, but got %#0.4x", e.DigitalSampleMask(), f.DigitalSampleMask())
					}
				}

				if f, ok := actual.(AnalogSampleMaskGetter); ok {
					e := tt.expected.(AnalogSampleMaskGetter)
					if f.AnalogSampleMask() != e.AnalogSampleMask() {
						t.Fatalf("Expected analog sample mask=%#0.2x, but got %#0.2x", e.AnalogSampleMask(), f.AnalogSampleMask())
					}
				}

				if f, ok := actual.(DigitalSamplesGetter); ok {
					e := tt.expected.(DigitalSamplesGetter)
					if f.DigitalSamples() != e.DigitalSamples() {
						t.Fatalf("Expected digital samples=%#0.4x, but got %#0.4x", e.DigitalSamples(), f.DigitalSamples())
					}
				}

				if f, ok := actual.(AnalogSampleGetter); ok {
					e := tt.expected.(AnalogSampleGetter)
					if f.AnalogSample() != e.AnalogSample() {
						t.Fatalf("Expected analog sample=%#0.4x, but got %#0.4x", e.AnalogSample(), f.AnalogSample())
					}
				}
			})
		}
	})
}

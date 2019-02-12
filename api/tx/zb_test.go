package tx

import (
	"testing"
)

var _ Frame = (*ZB)(nil)
var _ Addr64Setter = (*ZB)(nil)
var _ Addr16Setter = (*ZB)(nil)
var _ BroadcastRadiusSetter = (*ZB)(nil)
var _ OptionsSetter = (*ZB)(nil)
var _ DataSetter = (*ZB)(nil)

type zbTest struct {
	name     string
	input    *ZB
	expected []byte
}

var zbTests = []zbTest{
	{"ZB Defaults",
		NewZB(),
		[]byte{zbAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00}},
	{"ZB FrameID",
		NewZB(FrameID(1)),
		[]byte{zbAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00}},
	{"ZB Addresses",
		NewZB(Addr64(0x0001020304050607), Addr16(0x0102)),
		[]byte{zbAPIID, 0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01, 0x02, 0x00, 0x00}},
	{"ZB Broadcast Radius",
		NewZB(BroadcastRadius(0x10)),
		[]byte{zbAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x10, 0x00}},
	{"ZB Options",
		NewZB(Options(0x20)),
		[]byte{zbAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x20}},
	{"ZB Data",
		NewZB(Data([]byte{'h', 'e', 'l', 'l', 'o'})),
		[]byte{zbAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 'h', 'e', 'l', 'l', 'o'}},
}

func TestZB(t *testing.T) {
	t.Parallel()

	t.Run("ZB Test Suite", func(t *testing.T) {
		for _, tt := range zbTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				actual, err := tt.input.Bytes()
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				if len(actual) != len(tt.expected) {
					t.Fatalf("Expected ATRemote frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}
				for i, b := range actual {
					if b != tt.expected[i] {
						t.Fatalf("Expected 0x%02x, but got 0x%02x at index %d", b, actual[i], i)
					}
				}
			})
		}
	})
}

package tx

import (
	"fmt"
	"testing"
)

var _ Frame = (*ZBExplicit)(nil)
var _ Addr64Setter = (*ZBExplicit)(nil)
var _ Addr16Setter = (*ZBExplicit)(nil)
var _ SrcEPSetter = (*ZBExplicit)(nil)
var _ DstEPSetter = (*ZBExplicit)(nil)
var _ ClusterIDSetter = (*ZBExplicit)(nil)
var _ ProfileIDSetter = (*ZBExplicit)(nil)
var _ BroadcastRadiusSetter = (*ZBExplicit)(nil)
var _ OptionsSetter = (*ZBExplicit)(nil)
var _ DataSetter = (*ZBExplicit)(nil)

type zbExplicitTest struct {
	name     string
	input    *ZBExplicit
	expected []byte
}

var zbExplicitTests = []zbExplicitTest{
	{"ZB Explicit Defaults",
		NewZBExplicit(),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit FrameID",
		NewZBExplicit(FrameID(1)),
		[]byte{zbExplicitAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Addresses",
		NewZBExplicit(Addr64(0x0001020304050607), Addr16(0x0102)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Src EP",
		NewZBExplicit(SrcEP(0xaa)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Dst EP",
		NewZBExplicit(DstEP(0xaa)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Cluster ID",
		NewZBExplicit(ClusterID(0xaabb)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0xaa, 0xbb, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Profile ID",
		NewZBExplicit(ProfileID(0xaabb)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xbb, 0x00, 0x00}},
	{"ZB Explicit Broadcast Radius",
		NewZBExplicit(BroadcastRadius(0xaa)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0x00}},
	{"ZB Explicit Options",
		NewZBExplicit(Options(0xaa)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa}},
	{"ZB Explicit Data",
		NewZBExplicit(Data([]byte{'h', 'e', 'l', 'l', 'o'})),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 'h', 'e', 'l', 'l', 'o'}},
	{"ZB Explicit nil Data",
		NewZBExplicit(Data(nil)),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"ZB Explicit Empty Data",
		NewZBExplicit(Data([]byte{})),
		[]byte{zbExplicitAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

func TestZBExplicit(t *testing.T) {
	t.Parallel()

	t.Run("ZB Explicit Test Suite", func(t *testing.T) {
		for _, tt := range zbExplicitTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				actual, err := tt.input.Bytes()
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				if len(actual) != len(tt.expected) {
					var msg string
					for _, c := range tt.expected {
						msg = fmt.Sprintf("%s 0x%02X ", msg, c)
					}
					t.Logf("\texpected: %s", msg)

					msg = ""
					for _, c := range actual {
						msg = fmt.Sprintf("%s 0x%02X ", msg, c)
					}
					t.Logf("\t  actual: %s", msg)
					t.Fatalf("Expected ZB Explicit frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}
				for i, b := range actual {
					if b != byte(tt.expected[i]) {
						var msg string
						for _, c := range tt.expected {
							msg = fmt.Sprintf("%s 0x%02X ", msg, c)
						}
						t.Logf("\texpected: %s", msg)

						msg = ""
						for _, c := range actual {
							msg = fmt.Sprintf("%s 0x%02X ", msg, c)
						}
						t.Logf("\t  actual: %s", msg)
						t.Fatalf("Expected 0x%02x, but got 0x%02x at index %d", tt.expected[i], b, i)
					}
				}
			})
		}
	})
}

package tx

import (
	"testing"
)

var _ Frame = (*ATRemote)(nil)
var _ FrameIDSetter = (*ATRemote)(nil)
var _ Addr64Setter = (*ATRemote)(nil)
var _ Addr16Setter = (*ATRemote)(nil)
var _ OptionsSetter = (*ATRemote)(nil)
var _ CommandSetter = (*ATRemote)(nil)
var _ ParameterSetter = (*ATRemote)(nil)

type atRemoteTest struct {
	name     string
	input    *ATRemote
	expected []byte
}

var atRemoteTests = []atRemoteTest{
	{"AT Remote Defaults",NewATRemote(), []byte{atRemoteAPIID, 0, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote No Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(NI)), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote Nil Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(NI), Parameter(nil)), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote Empty Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(NI), Parameter([]byte{})), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote With Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(NI), Parameter([]byte{1})), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I', 0x01}},
}

func TestATRemote(t *testing.T) {
	t.Parallel()

	t.Run("AT Remote Test Suite", func(t *testing.T) {
		for _, tt := range atRemoteTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				actual, err := tt.input.Bytes()
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				if len(actual) != len(tt.expected) {
					t.Fatalf("Expected frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
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

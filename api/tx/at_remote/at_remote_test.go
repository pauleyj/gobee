package at_remote

import (
	"testing"
	"github.com/pauleyj/gobee/api/tx"
)

var _ tx.Frame = (*ATRemote)(nil)

type atTest struct {
	name     string
	input    *ATRemote
	expected []byte
}

var atTests = []atTest{
	{"AT Remote No Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(tx.NI)), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote Nil Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(tx.NI), Parameter(nil)), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote Empty Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(tx.NI), Parameter([]byte{})), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I'}},
	{"AT Remote With Param", NewATRemote(FrameID(0x01), Addr64(0x000000000000FFFF), Addr16(0xFFFE), Options(0), Command(tx.NI), Parameter([]byte{1})), []byte{atRemoteAPIID, 1, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xff, 0xff, 0xff, 0xfe, 0x00, 'N', 'I', 0x01}},
}

func TestAT(t *testing.T) {
	t.Parallel()

	t.Run("AT Queue Test Suite", func(t *testing.T) {
		t.Parallel()
		for _, tt := range atTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				actual, err := tt.input.Bytes()
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if len(actual) != len(tt.expected) {
					t.Errorf("Expected ATRemote frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}
				for i, b := range actual {
					if b != tt.expected[i] {
						t.Errorf("Expected 0x%02x, but got 0x%02x at index %d", b, actual[i], i)
					}
				}
			})
		}
	})
}

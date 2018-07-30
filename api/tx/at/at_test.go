package at

import (
	"testing"
	"github.com/pauleyj/gobee/api/tx"
)

var _ tx.Frame = (*AT)(nil)

type atTest struct {
	name     string
	input    *AT
	expected []byte
}

var atTests = []atTest{
	{"AT No Param", NewAT(FrameID(0x01), Command(tx.NI)), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT Nil Param", NewAT(FrameID(0x01), Command(tx.NI), Parameter(nil)), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT Empty Param", NewAT(FrameID(0x01), Command(tx.NI), Parameter([]byte{})), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT With Param", NewAT(FrameID(0x01), Command(tx.NI), Parameter([]byte{0x00})), []byte{atAPIID, 1, 'N', 'I', 0}},
}

func TestAT(t *testing.T) {
	t.Parallel()

	t.Run("AT Test Suite", func(t *testing.T) {
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
					t.Errorf("Expected AT frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}
				for i, b := range actual {
					if b != tt.expected[i] {
						t.Errorf("Expected 0x%02x, but got 0x%02x", b, actual[i])
					}
				}
			})
		}
	})
}
package tx

import (
	"testing"
)

var _ Frame = (*AT)(nil)
var _ FrameIDSetter = (*AT)(nil)
var _ CommandSetter = (*AT)(nil)
var _ ParameterSetter = (*AT)(nil)

type atTest struct {
	name     string
	input    *AT
	expected []byte
}

var atTests = []atTest{
	{"AT Defaults", NewAT(), []byte{atAPIID, 0, 'N', 'I'}},
	{"AT No Param", NewAT(FrameID(0x01), Command(NI)), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT Nil Param", NewAT(FrameID(0x01), Command(NI), Parameter(nil)), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT Empty Param", NewAT(FrameID(0x01), Command(NI), Parameter([]byte{})), []byte{atAPIID, 1, 'N', 'I'}},
	{"AT With Param", NewAT(FrameID(0x01), Command(NI), Parameter([]byte{0x00})), []byte{atAPIID, 1, 'N', 'I', 0}},
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
					t.Fatalf("Expected no error, but got: %v", err)
				}
				if len(actual) != len(tt.expected) {
					t.Fatalf("Expected AT frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}
				for i, b := range actual {
					if b != tt.expected[i] {
						t.Fatalf("Expected 0x%02x, but got 0x%02x", b, actual[i])
					}
				}
			})
		}
	})
}

package tx

import (
	"testing"
)

var _ Frame = (*ATQueue)(nil)
var _ FrameIDSetter = (*ATQueue)(nil)
var _ CommandSetter = (*ATQueue)(nil)
var _ ParameterSetter = (*ATQueue)(nil)

type atQueueTest struct {
	name     string
	input    *ATQueue
	expected []byte
}

var atQueueTests = []atQueueTest{
	{"AT Queue Default", NewATQueue(), []byte{atQueueAPIID, 0, 0, 0}},
	{"AT Queue No Param", NewATQueue(FrameID(0x01), Command(NI)), []byte{atQueueAPIID, 1, 'N', 'I'}},
	{"AT Queue Nil Param", NewATQueue(FrameID(0x01), Command(NI), Parameter(nil)), []byte{atQueueAPIID, 1, 'N', 'I'}},
	{"AT Queue Empty Param", NewATQueue(FrameID(0x01), Command(NI), Parameter([]byte{})), []byte{atQueueAPIID, 1, 'N', 'I'}},
	{"AT Queue With Param", NewATQueue(FrameID(0x01), Command(NI), Parameter([]byte{0x00})), []byte{atQueueAPIID, 1, 'N', 'I', 0}},
}

func TestATQueue(t *testing.T) {
	t.Parallel()

	t.Run("AT Queue Test Suite", func(t *testing.T) {
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
					t.Fatalf("Expected frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
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
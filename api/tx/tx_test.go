package tx

import (
	"fmt"
	"github.com/pauleyj/gobee/api"
	"testing"
)

type dummyFrame struct {
	data []byte
}

func (f *dummyFrame) Bytes() ([]byte, error) {
	return f.data, nil
}

type txFrameTest struct {
	name     string
	input    Frame
	f        *APIFrame
	expected []byte
}

var txFrameTests = []txFrameTest{
	{"API Frame Default",
		&dummyFrame{data: []byte{0x08, 0x01, 'N', 'J'}},
		New(),
		[]byte{0x7E, 0x00, 0x04, 0x08, 0x01, 'N', 'J', 0x5e}},
	{"API Frame nil Options",
		&dummyFrame{data: []byte{0x08, 0x01, 'N', 'J'}},
		New(nil),
		[]byte{0x7E, 0x00, 0x04, 0x08, 0x01, 'N', 'J', 0x5e}},
	{"API Frame No Escape",
		&dummyFrame{data: []byte{0x08, 0x01, 'N', 'J'}},
		New(api.APIEscapeMode(api.EscapeModeInactive)),
		[]byte{0x7E, 0x00, 0x04, 0x08, 0x01, 'N', 'J', 0x5e}},
	{"API Frame With Escape",
		&dummyFrame{[]byte{0x23, 0x7E, 0x7D, 0x11, 0x13}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{0x7E, 0x00, 0x05, 0x23, 0x7D, 0x5E, 0x7D, 0x5D, 0x7D, 0x31, 0x7D, 0x33, 0xBD}},
	{"API Frame With Escape",
		&dummyFrame{[]byte{0x23, 0x11}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{0x7E, 0x00, 0x02, 0x23, 0x7D, 0x31, 0xcb}},
	{"API Frame with Checksum Escaped",
		&dummyFrame{[]byte{
			0x10, 0x01, 0x00, 0x13,
			0xA2, 0x00, 0x40, 0x0A,
			0x01, 0x27, 0xFF, 0xFE,
			0x00, 0x00, 0x54, 0x78,
			0x44, 0x61, 0x74, 0x61,
			0x30, 0x41}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{
			0x7E, 0x00, 0x16, 0x10,
			0x01, 0x00, 0x7D, 0x33,
			0xA2, 0x00, 0x40, 0x0A,
			0x01, 0x27, 0xFF, 0xFE,
			0x00, 0x00, 0x54, 0x78,
			0x44, 0x61, 0x74, 0x61,
			0x30, 0x41, 0x7D, 0x33},
	},
	{"API Frame With Payload Escape",
		&dummyFrame{[]byte{
			0x10,
			0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x13}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{
			0x7e, 0x00, 0x0f, 0x10,
			0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x7d, 0x33, 0xdb},
	},
	{"API Frame With Payload Escape",
		&dummyFrame{[]byte{
			0x10, 0x01, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0xff, 0xff, 0x5f, 0xd6,
			0x00, 0x00, 0x13}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{
			0x7e, 0x00, 0x0f, 0x10,
			0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0xff,
			0xff, 0x5f, 0xd6, 0x00,
			0x00, 0x7d, 0x33, 0xa8},
	},
	{"API Frame With Length Escape",
		&dummyFrame{[]byte{
			0x10, 0x04, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0xff, 0xff, 0x5f, 0xd6,
			0x00, 0x00, 0x61, 0x73,
			0x64}},
		New(api.APIEscapeMode(api.EscapeModeActive)),
		[]byte{
			0x7e, 0x00, 0x7d, 0x31,
			0x10, 0x04, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0xff, 0xff, 0x5f, 0xd6,
			0x00, 0x00, 0x61, 0x73,
			0x64, 0x80},
	},
	//
}

func TestTXAPIFrame(t *testing.T) {
	t.Parallel()

	t.Run("TX API Suite", func(t *testing.T) {
		for _, tt := range txFrameTests {
			tt := tt

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				actual, err := tt.f.Bytes(tt.input)
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				//msg := fmt.Sprintf("%10s:", "a")
				//for _, c := range actual {
				//	msg = fmt.Sprintf("%s %#0.2x", msg, c)
				//}
				//t.Log(msg)
				//
				//msg = fmt.Sprintf("%10s:", "e")
				//for _, c := range actual {
				//	msg = fmt.Sprintf("%s %#0.2x", msg, c)
				//}
				//t.Log(msg)

				if len(actual) != len(tt.expected) {
					msg := fmt.Sprintf("%10s:", "a")
					for _, c := range actual {
						msg = fmt.Sprintf("%s %#0.2x", msg, c)
					}
					t.Log(msg)

					msg = fmt.Sprintf("%10s:", "e")
					for _, c := range tt.expected {
						msg = fmt.Sprintf("%s %#0.2x", msg, c)
					}
					t.Log(msg)
					t.Fatalf("Expected API frame to be %d bytes in length, got: %d", len(tt.expected), len(actual))
				}

				for i, b := range actual {
					if b != tt.expected[i] {
						msg := fmt.Sprintf("%10s:", "a")
						for _, c := range actual {
							msg = fmt.Sprintf("%s %#0.2x", msg, c)
						}
						t.Log(msg)

						msg = fmt.Sprintf("%10s:", "e")
						for _, c := range tt.expected {
							msg = fmt.Sprintf("%s %#0.2x", msg, c)
						}
						t.Log(msg)
						t.Fatalf("Expected 0x%02x, but got 0x%02x at index %d", tt.expected[i], b, i)
					}
				}
			})
		}
	})
}

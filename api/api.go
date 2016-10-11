package api

import "errors"

// FrameDelimiter start API frame delimiter, requires escaping in mode 2
const FrameDelimiter byte = 0x7E

// dataLengthBytes Number of data length bytes
const FrameLengthByteCount uint16 = 2

// validChecksum API frame valid checksum
const ValidChecksum byte = 0xFF

// ESC the escape character
const ESC byte = 0x7D

// xon XON character, requires escaping in mode 2
const xon byte = 0x11

// xoff XOFF character, requires escaping in mode 2
const xoff byte = 0x13

// EscChar the character used to escape charters needing escaping
const ESCChar byte = 0x20

var (
	escapeSet             = [...]byte{FrameDelimiter, ESC, xon, xoff}
	ErrChecksumValidation = errors.New("Frame failed checksum validation")
	ErrFrameDelimiter     = errors.New("Expected frame delimiter")
	ErrInvalidAPIMode     = errors.New("Invalid API mode")
)

// State the API frame state type
type State int

const (
	FrameStart    = State(iota)
	FrameLength   = State(iota)
	APIID         = State(iota)
	FrameData     = State(iota)
	FrameChecksum = State(iota)
)

// APIEscapeMode defines the XBee API escape mode type
type APIEscapeMode byte

const (
	EscapeModeInactive = APIEscapeMode(1)
	EscapeModeActive   = APIEscapeMode(2)
)

func ShouldEscape(c byte) bool {
	return include(escapeSet[:], c)
}

func Escape(c byte) byte {
	return c ^ ESCChar
}

// index returns the first index of the target byte t, or -1 if no match is found
func index(vc []byte, c byte) int {
	for i, v := range vc {
		if v == c {
			return i
		}
	}
	return -1
}

// include returns true if the target byte t is in the slice.
func include(vc []byte, c byte) bool {
	return index(vc, c) >= 0
}

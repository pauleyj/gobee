package api

import "errors"

// FrameDelimiter start API frame delimiter, requires escaping in mode 2
const FrameDelimiter byte = 0x7E

// FrameLengthByteCount Number of data length bytes
const FrameLengthByteCount uint16 = 2

// ValidChecksum API frame valid checksum
const ValidChecksum byte = 0xFF

// ESC the escape character
const ESC byte = 0x7D

// xon XON character, requires escaping in mode 2
const xon byte = 0x11

// xoff XOFF character, requires escaping in mode 2
const xoff byte = 0x13

// ESCChar the character used to escape charters needing escaping
const ESCChar byte = 0x20

var (
	escapeSet = map[byte]struct{}{
		FrameDelimiter: {},
		ESC:            {},
		xon:            {},
		xoff:           {},
	}
	// ErrChecksumValidation frame failed checksum validation
	ErrChecksumValidation = errors.New("Frame failed checksum validation")
	// ErrFrameDelimiter expecting frame start delimiter
	ErrFrameDelimiter = errors.New("Expected frame delimiter")
	// ErrInvalidAPIEscapeMode invalid API escape mode
	ErrInvalidAPIEscapeMode = errors.New("Invalid API escape mode")
)

// State the API frame state type
type State int

// Frame states
const (
	FrameStart    = State(iota)
	FrameLength   = State(iota)
	APIID         = State(iota)
	FrameData     = State(iota)
	FrameChecksum = State(iota)
)

// EscapeMode defines the XBee API escape mode type
type EscapeMode byte

// Escape modes
const (
	EscapeModeInactive = EscapeMode(1)
	EscapeModeActive   = EscapeMode(2)
)

// ShouldEscape should this byte be escaped
func ShouldEscape(c byte) bool {
	_, ok := escapeSet[c]
	return ok
}

// Escape escape this byte
func Escape(c byte) byte {
	return c ^ ESCChar
}

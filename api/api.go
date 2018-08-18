package api

import "errors"

// BroadcastAddr64 64-bit broadcast address
const BroadcastAddr64 uint64 = 0x000000000000FFFF

// BroadcastAddr16 16-bit broadcast address
const BroadcastAddr16 uint16 = 0xFFFE

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
	ErrChecksumValidation = errors.New("checksum validation error")
	// ErrFrameDelimiter expecting frame start delimiter
	ErrFrameDelimiter = errors.New("expected frame delimiter")
	// ErrInvalidAPIEscapeMode invalid API escape mode
	ErrInvalidAPIEscapeMode = errors.New("invalid API escape mode")
)

// State the API frame state type
type State int

// Frame internal states
const (
	FrameStart    = State(iota)
	FrameLength   = State(iota)
	APIID         = State(iota)
	FrameData     = State(iota)
	FrameChecksum = State(iota)
)

// EscapeMode defines the XBee API escape mode type
type EscapeMode byte

// Escape valid escape modes
const (
	EscapeModeInactive = EscapeMode(1)
	EscapeModeActive   = EscapeMode(2)
)

// APIEscapeModeSetter interface for APIEscapeMode setters
type APIEscapeModeSetter interface {
	SetAPIEscapeMode(EscapeMode)
}

// APIEscapeMode options helper function for APIEscapeModeSetter
func APIEscapeMode(mode EscapeMode) func(interface{}) {
	return func(i interface{}) {
		if t, ok := i.(APIEscapeModeSetter); ok {
			t.SetAPIEscapeMode(mode)
		}
	}
}

// ShouldEscape should this byte be escaped
func ShouldEscape(c byte) bool {
	_, ok := escapeSet[c]
	return ok
}

// Escape escape this byte
func Escape(c byte) byte {
	return c ^ ESCChar
}

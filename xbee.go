package gobee

import (
	"bytes"
	"errors"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
)

// frameDelimiter start API frame delimiter, requires escaping in mode 2
const frameDelimiter byte = 0x7E

// dataLengthBytes Number of data length bytes
const dataLengthBytes uint8 = 2

// validChecksum API frame valid checksum
const validChecksum byte = 0xFF

// esc escape character
const esc byte = 0x7D

// xon XON character, requires escaping in mode 2
const xon byte = 0x11

// xoff XOFF character, requires escaping in mode 2
const xoff byte = 0x13

// escChar the character used to escape charters needing escaping
const escChar = 0x20

// BroadcastAddr64 64-bit broadcast address
const BroadcastAddr64 uint64 = 0x000000000000FFFF

// BriadcastAddr16 16-bit broadcast address
const BriadcastAddr16 uint16 = 0xFFFE

var (
	escapeSet             = [...]byte{frameDelimiter, esc, xon, xoff}
	errChecksumValidation = errors.New("Frame failed checksum validation")
	errFrameDelimiter     = errors.New("Expected frame delimiter")
	errInvalidAPIMode     = errors.New("Invalid API mode")
)

// apiState the API state type
type apiState int

const (
	frameStart    = apiState(iota)
	frameLength   = apiState(iota)
	apiID         = apiState(iota)
	frameData     = apiState(iota)
	frameChecksum = apiState(iota)
)

// XBeeTransmitter used to transmit API frame bytes to serial communications port
type XBeeTransmitter interface {
	Transmit([]byte) (int, error)
}

// XBeeReceiver used to report frames received by the XBee
type XBeeReceiver interface {
	Receive(rx.Frame) error
}

// XBee all the things
type XBee struct {
	transmitter              XBeeTransmitter
	receiver                 XBeeReceiver
	apiMode                  byte
	escapeNext               bool
	rxState                  apiState
	rxFrameDataSize          uint16
	rxFrameChecksum          uint8
	rxFrameDataSizeByteIndex uint8
	rxFrameDataIndex         uint16
	rxFrame                  rx.Frame
}

// NewXBee constructor of XBee's
func New(transmitter XBeeTransmitter, receiver XBeeReceiver) *XBee {
	return &XBee{
		transmitter:              transmitter,
		receiver:                 receiver,
		apiMode:                  1,
		escapeNext:               false,
		rxState:                  frameStart,
		rxFrameDataSize:          0,
		rxFrameChecksum:          0,
		rxFrameDataSizeByteIndex: 0,
		rxFrameDataIndex:         0,
		rxFrame:                  nil,
	}
}

// RX bytes received from the serial communications port are sent here
func (x *XBee) RX(b byte) error {
	var err error

	if x.isAPIEscapeModeEnabled() {
		if x.rxState != frameStart && b == esc && !x.escapeNext {
			x.escapeNext = true
			return nil
		}

		if x.escapeNext {
			x.escapeNext = false
			b = escape(b)
		}
	}

	switch x.rxState {
	case frameLength:
		err = x.apiStateDataLength(b)
	case apiID:
		err = x.apiStateAPIID(b)
	case frameData:
		err = x.apiStateFrameData(b)
	case frameChecksum:
		err = x.apiStateChecksum(b)
		if err == nil {
			x.receiver.Receive(x.rxFrame)
		}
	default:
		err = x.apiStateWaitFrameDelimiter(b)
	}
	return err
}

// TX transmit a frame to the XBee, forms an appropriate API frame for the frame being sent,
// uses the XBeeTransmitter to send the API frame bytes to the serial communications port
func (x *XBee) TX(frame tx.Frame) (int, error) {
	f, err := frame.Bytes()
	if err != nil {
		return 0, err
	}

	var b bytes.Buffer
	var checksum byte
	l := len(f)

	b.WriteByte(frameDelimiter)

	lh := byte(l >> 8)
	ll := byte(l & 0x00FF)
	if x.isAPIEscapeModeEnabled() {
		if shouldEscape(lh) {
			lh = escape(lh)
			b.WriteByte(esc)
		}
		if shouldEscape(ll) {
			ll = escape(ll)
			b.WriteByte(esc)
		}
	}
	b.WriteByte(lh)
	b.WriteByte(ll)

	for _, i := range f {
		// checksum is calculated on pre escaped value
		checksum += i

		if x.isAPIEscapeModeEnabled() && shouldEscape(i) {
			i = escape(i)
			b.WriteByte(esc)
		}

		b.WriteByte(i)
	}

	checksum = validChecksum - checksum

	// checksum is escaped if needed
	if x.apiMode == 2 && shouldEscape(checksum) {
		checksum = escape(checksum)
		b.WriteByte(esc)
	}
	b.WriteByte(checksum)

	return x.transmitter.Transmit(b.Bytes())
}

// SetAPIMode sets the API mode so goobe knows to escape or not
func (x *XBee) SetAPIMode(mode byte) error {
	if mode != 1 && mode != 2 {
		return errInvalidAPIMode
	}
	x.apiMode = mode
	return nil
}

func shouldEscape(b byte) bool {
	for _, v := range escapeSet {
		if b == v {
			return true
		}
	}
	return false
}

func escape(b byte) byte {
	return (b ^ escChar)
}

func (x *XBee) isAPIEscapeModeEnabled() bool {
	return x.apiMode == 0x02
}

func (x *XBee) apiStateWaitFrameDelimiter(b byte) error {
	if frameDelimiter == b {
		x.rxFrameDataSize = 0
		x.rxFrameChecksum = 0
		x.rxFrameDataSizeByteIndex = 0
		x.rxFrameDataIndex = 0
		x.rxFrame = nil
		x.rxState = frameLength

		return nil
	}

	x.rxState = frameStart
	return errFrameDelimiter
}

func (x *XBee) apiStateDataLength(b byte) error {
	x.rxFrameDataSize += uint16(b << (1 - x.rxFrameDataSizeByteIndex))
	x.rxFrameDataSizeByteIndex++

	if x.rxFrameDataSizeByteIndex == dataLengthBytes {
		x.rxState = apiID
	}

	return nil
}

func (x *XBee) apiStateAPIID(b byte) error {
	var err error

	x.rxFrame, err = rx.NewFrameForAPIID(b)
	if err != nil {
		x.rxState = frameStart
		return err
	}
	x.rxFrameChecksum += b
	x.rxFrameDataIndex++
	x.rxState = frameData

	return nil
}

func (x *XBee) apiStateFrameData(b byte) error {
	x.rxFrame.RX(b)
	x.rxFrameChecksum += b
	x.rxFrameDataIndex++

	if x.rxFrameDataIndex == x.rxFrameDataSize {
		x.rxState = frameChecksum
	}

	return nil
}

func (x *XBee) apiStateChecksum(b byte) error {
	x.rxState = frameStart
	x.rxFrameChecksum += b

	if validChecksum != x.rxFrameChecksum {
		return errChecksumValidation
	}

	return nil
}

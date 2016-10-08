package gobee

import (
	"bytes"
	"errors"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
	"encoding/binary"
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
const BroadcastAddr16 uint16 = 0xFFFE

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

// APIEscapeMode defines the XBee API escape mode type
type APIEscapeMode byte

const (
	EscapeModeInactive = APIEscapeMode(1)
	EscapeModeActive = APIEscapeMode(2)
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
	apiMode                  APIEscapeMode
	escapeNext               bool
	rxState                  apiState
	rxFrameDataSize          uint16
	rxFrameChecksum          uint8
	rxFrameDataSizeByteIndex uint8
	rxFrameDataIndex         uint16
	rxFrame                  rx.Frame
}

// New constructor of XBee's
func New(transmitter XBeeTransmitter, receiver XBeeReceiver) *XBee {
	return &XBee{
		transmitter:              transmitter,
		receiver:                 receiver,
		apiMode:                  EscapeModeInactive,
		escapeNext:               false,
		rxState:                  frameStart,
		rxFrameDataSize:          0,
		rxFrameChecksum:          0,
		rxFrameDataSizeByteIndex: 0,
		rxFrameDataIndex:         0,
		rxFrame:                  nil,
	}
}

func NewWithEscapeMode(transmitter XBeeTransmitter, receiver XBeeReceiver, mode APIEscapeMode) *XBee {
	return &XBee{
		transmitter:              transmitter,
		receiver:                 receiver,
		apiMode:                  mode,
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

	return x.handleRX(b)
}

func (x *XBee) handleRX(b byte) error {
	var err error
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
	b.WriteByte(frameDelimiter)
	b.Write(uint16ToBytes(uint16(len(f))))
	b.Write(f)
	b.WriteByte(checksum(f))

	var i int
	if x.apiMode == EscapeModeInactive {
		i, err = x.transmitter.Transmit(b.Bytes())
	} else {
		i, err = x.transmitter.Transmit(escapeTXBuffer(b.Bytes()))
	}

	return i, err
}

// SetAPIMode sets the API mode so goobe knows to escape or not
func (x *XBee) SetAPIMode(mode APIEscapeMode) error {
	if mode != EscapeModeInactive && mode != EscapeModeActive {
		return errInvalidAPIMode
	}
	x.apiMode = mode
	return nil
}

func escapeTXBuffer(b []byte) []byte {
	var e bytes.Buffer
	for i, c := range b {
		if i != 0 && shouldEscape(c) {
			c = escape(c)
			e.WriteByte(escChar)
		}
		e.WriteByte(c)
	}

	return e.Bytes()
}

func checksum(buff []byte) byte {
	var c byte
	for _, b := range buff {
		c += b
	}
	return validChecksum - c
}

func escape(b byte) byte {
	return b ^ escChar
}

func shouldEscape(b byte) bool {
	for _, v := range escapeSet {
		if b == v {
			return true
		}
	}
	return false
}

func (x *XBee) isAPIEscapeModeEnabled() bool {
	return x.apiMode == EscapeModeActive
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

func uint16ToBytes(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

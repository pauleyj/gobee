package gobee

import (
	"bytes"
	"errors"
	"github.com/pauleyj/gobee/rx"
	"github.com/pauleyj/gobee/tx"
	"io"
)

/*
 * Frame Constants
 */
const FRAME_DELIMITER byte = 0x7E
const XBEE_NUMBER_OF_SIZE_BYTES uint8 = 2
const XBEE_VALID_FRAME_CHECKSUM byte = 0xFF
const ESC byte = 0x7D
const XON byte = 0x11
const XOFF byte = 0x13
const ESC_CHAR = 0x20

/*
 * Address constants
 */
const BROADCAST_ADDR_64 uint64 = 0x000000000000FFFF
const BROADCAST_ADDR_16 uint16 = 0xFFFE

var (
	escape_set = [...]byte{FRAME_DELIMITER, ESC, XON, XOFF}
	errChecksumValidation = errors.New("Frame failed checksum validation")
	errInvalidApiMode = errors.New("Invalid API mode")
)

const (
	STATE_DATA_FRAME_START = ApiState(iota)
	STATE_DATA_FRAME_LENGTH = ApiState(iota)
	STATE_DATA_FRAME_API_ID = ApiState(iota)
	STATE_DATA_FRAME_DATA = ApiState(iota)
	STATE_DATA_FRAME_CHECKSUM = ApiState(iota)
)

type ApiState int

type XBeeTransmitter interface {
	io.Writer
}

type XBeeReceiver interface {
	RxFrameReceiver(rx.RxFrame) error
}

type XBee struct {
	transmitter              XBeeTransmitter
	receiver                 XBeeReceiver
	apiMode                  byte
	escapeNext               bool
	rxState                  ApiState
	rxFrameDataSize          uint16
	rxFrameChecksum          uint8
	rxFrameDataSizeByteIndex uint8
	rxFrameDataIndex         uint16
	rxFrame                  rx.RxFrame
}

func NewXBee(transmitter XBeeTransmitter, receiver XBeeReceiver) *XBee {
	return &XBee{
		transmitter:              transmitter,
		receiver:                 receiver,
		apiMode: 		  1,
		escapeNext:               false,
		rxState:                  STATE_DATA_FRAME_START,
		rxFrameDataSize:          0,
		rxFrameChecksum:          0,
		rxFrameDataSizeByteIndex: 0,
		rxFrameDataIndex:         0,
		rxFrame:                  nil,
	}
}

func (x *XBee) RX(b byte) error {
	var err error

	if x.isApiEscapeModeEnabled() {
		if x.rxState != STATE_DATA_FRAME_START && b == ESC && !x.escapeNext {
			x.escapeNext = true
			return nil
		}

		if x.escapeNext {
			x.escapeNext = false
			b = escape(b)
		}
	}

	switch x.rxState {
	case STATE_DATA_FRAME_LENGTH:
		err = x.apiStateDataLength(b)
	case STATE_DATA_FRAME_API_ID:
		err = x.apiStateApiId(b)
	case STATE_DATA_FRAME_DATA:
		err = x.apiStateFrameData(b)
	case STATE_DATA_FRAME_CHECKSUM:
		err = x.apiStateChecksum(b)
		if err == nil {
			x.receiver.RxFrameReceiver(x.rxFrame)
		}
	default:
		err = x.apiStateWaitFrameDelimiter(b)
	}
	return err
}


func (x *XBee) TX(frame tx.TxFrame) (int, error) {
	f, err := frame.Bytes()
	if err != nil {
		return 0, err
	}

	var b bytes.Buffer
	var checksum byte = 0
	l := len(f)

	b.WriteByte(FRAME_DELIMITER)

	lh := byte(l >> 8)
	ll := byte(l & 0x00FF)
	if x.isApiEscapeModeEnabled() {
		if shouldEscape(lh) {
			lh = escape(lh)
			b.WriteByte(ESC)
		}
		if shouldEscape(ll) {
			ll = escape(ll)
			b.WriteByte(ESC)
		}
	}
	b.WriteByte(lh)
	b.WriteByte(ll)

	for _, i := range f {
		// checksum is calculated on pre escaped value
		checksum += i

		if x.isApiEscapeModeEnabled() && shouldEscape(i) {
			i = escape(i)
			b.WriteByte(ESC)
		}

		b.WriteByte(i)
	}

	checksum = XBEE_VALID_FRAME_CHECKSUM - checksum

	// checksum is escaped if needed
	if x.apiMode == 2 && shouldEscape(checksum) {
		checksum = escape(checksum)
		b.WriteByte(ESC)
	}
	b.WriteByte(checksum)

	return x.transmitter.Write(b.Bytes())
}

func (x *XBee) SetApiMode(mode byte) error {
	if mode != 1 && mode != 2 {
		return errInvalidApiMode
	}
	x.apiMode = mode
	return nil
}

func shouldEscape(b byte) bool {
	for _, v := range escape_set {
		if b == v {
			return true
		}
	}
	return false
}

func escape(b byte) byte {
	return (b ^ ESC_CHAR)
}

func (x *XBee) isApiEscapeModeEnabled() bool {
	return x.apiMode == 0x02
}


func (x *XBee) apiStateWaitFrameDelimiter(b byte) error {
	if FRAME_DELIMITER == b {
		x.rxFrameDataSize = 0
		x.rxFrameChecksum = 0
		x.rxFrameDataSizeByteIndex = 0
		x.rxFrameDataIndex = 0
		x.rxFrame = nil
		x.rxState = STATE_DATA_FRAME_LENGTH

		return nil
	}

	x.rxState = STATE_DATA_FRAME_START

	return nil
}

func (x *XBee) apiStateDataLength(b byte) error {
	x.rxFrameDataSize += uint16(b << (1 - x.rxFrameDataSizeByteIndex))
	x.rxFrameDataSizeByteIndex++

	if x.rxFrameDataSizeByteIndex == XBEE_NUMBER_OF_SIZE_BYTES {
		x.rxState = STATE_DATA_FRAME_API_ID
	}

	return nil
}

func (x *XBee) apiStateApiId(b byte) error {
	var err error

	x.rxFrame, err = rx.NewRxFrameForApiId(b)
	if err != nil {
		x.rxState = STATE_DATA_FRAME_START
		return err
	}
	x.rxFrameChecksum += b
	x.rxFrameDataIndex++
	x.rxState = STATE_DATA_FRAME_DATA

	return nil
}

func (x *XBee) apiStateFrameData(b byte) error {
	x.rxFrame.RX(b)
	x.rxFrameChecksum += b
	x.rxFrameDataIndex++

	if x.rxFrameDataIndex == x.rxFrameDataSize {
		x.rxState = STATE_DATA_FRAME_CHECKSUM
	}

	return nil
}

func (x *XBee) apiStateChecksum(b byte) error {
	x.rxState = STATE_DATA_FRAME_START
	x.rxFrameChecksum += b

	if XBEE_VALID_FRAME_CHECKSUM != x.rxFrameChecksum {
		return errChecksumValidation
	}

	return nil
}

package rx

import "errors"

// frameFactory private type that defines a function returning an RxFrame
type FrameFactory func() RxFrame

var (
	errUnknownFrameApiId = errors.New("Unknown frame API ID")

	rxFrameFactory = make(map[byte]FrameFactory)
)

func init() {
	// AT command response
	rxFrameFactory[XBEE_API_ID_AT_COMMAND_RESPONSE] = newATCommandResponse
	rxFrameFactory[XBEE_API_ID_RX_ZB] = newRX_ZB
}

// NewRxFrameForApiId creates an appropriate RxFrame for the given API ID
func NewRxFrameForApiId(id byte) (RxFrame, error) {
	if f, ok := rxFrameFactory[id]; ok {
		return f(), nil
	}
	return nil, errUnknownFrameApiId
}

func AddApiFactoryForId(id byte, factory FrameFactory ) {
	if _, ok := rxFrameFactory[id]; !ok {
		rxFrameFactory[id] = factory
	}
}
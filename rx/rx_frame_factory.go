package rx

import "errors"

// defines a function returning an RxFrame
type FrameFactory func() RxFrame

var (
	errUnknownFrameApiId = errors.New("Unknown frame API ID")
	errFrameIdExists = errors.New("Frame factory for API ID already exists")
	rxFrameFactory map[byte]FrameFactory
)

func init() {
	rxFrameFactory = make(map[byte]FrameFactory)
	// AT command response
	rxFrameFactory[XBEE_API_ID_RX_AT] = newAT
	rxFrameFactory[XBEE_API_ID_RX_ZB] = newZB
	rxFrameFactory[XBEE_API_TX_STATUS] = newTX_STATUS
	rxFrameFactory[XBEE_API_ID_RX_ZB_EXPLICIT] = newZB_EXPLICIT
	rxFrameFactory[XBEE_API_ID_RX_AT_REMOTE] = newAT_REMOTE
}

// NewRxFrameForApiId creates an appropriate RxFrame for the given API ID
func NewRxFrameForApiId(id byte) (RxFrame, error) {
	if f, ok := rxFrameFactory[id]; ok {
		return f(), nil
	}
	return nil, errUnknownFrameApiId
}

func AddApiFactoryForId(id byte, factory FrameFactory ) error {
	if _, ok := rxFrameFactory[id]; !ok {
		rxFrameFactory[id] = factory
	} else {
		return errFrameIdExists
	}

	return nil
}
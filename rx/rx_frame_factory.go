package rx

import "errors"

// defines a function returning an RxFrame
type FrameFactory func() RxFrame

var (
	errUnknownFrameApiId = errors.New("Unknown frame API ID")
	errFrameIdExists = errors.New("Frame factory for API ID already exists")
	rxFrameFactory = make(map[byte]FrameFactory)
)

func init() {
	// AT command response
	rxFrameFactory[XBEE_API_ID_RX_AT] = newAT
	rxFrameFactory[XBEE_API_ID_RX_ZB] = newZB
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
package rx

import "errors"

// FrameFactory defines a function returning an RX Frame
type FrameFactory func() Frame

var (
	errUnknownFrameAPIID = errors.New("Unknown frame API ID")
	errFrameAPIIDExists  = errors.New("Frame factory for API ID already exists")
	rxFrameFactory       map[byte]FrameFactory
)

func init() {
	rxFrameFactory = make(map[byte]FrameFactory)

	rxFrameFactory[atAPIID] = newAT
	rxFrameFactory[zbAPIID] = newZB
	rxFrameFactory[txStatusAPIID] = newTXStatus
	rxFrameFactory[zbExplicitAPIID] = newZBExplicit
	rxFrameFactory[atRemoteAPIID] = newATRemote
	rxFrameFactory[modemStatusAPIID] = newModemStatus
}

// NewFrameForAPIID creates an appropriate RxFrame for the given API ID
func NewFrameForAPIID(id byte) (Frame, error) {
	if f, ok := rxFrameFactory[id]; ok {
		return f(), nil
	}
	return nil, errUnknownFrameAPIID
}

// AddFactoryForAPIID add frame by ID so factory can produce
func AddFactoryForAPIID(id byte, factory FrameFactory) error {
	if _, ok := rxFrameFactory[id]; !ok {
		rxFrameFactory[id] = factory
	} else {
		return errFrameAPIIDExists
	}

	return nil
}

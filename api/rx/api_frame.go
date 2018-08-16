package rx

import (
	"github.com/pauleyj/gobee/api"
)

func NewAPIFrame(options ...func(interface{})) *APIFrame {
	f := &APIFrame{}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// APIFrame defines an RX API frame
type APIFrame struct {
	Mode  api.EscapeMode
	state state
	frame Frame
}

func (f *APIFrame) SetAPIEscapeMode(mode api.EscapeMode) {
	f.Mode = mode
}

// RX receive byte
func (f *APIFrame) RX(c byte) (Frame, error) {
	if f.Mode == api.EscapeModeActive {

		if f.shouldEscapeNext(c) {
			return nil, nil
		}
		if f.state.escape {
			c = f.escape(c)
		}
	}

	return f.processRX(c)
}

func (f *APIFrame) processRX(c byte) (Frame, error) {
	switch f.state.state {
	case api.FrameLength:
		return nil, f.handleStateLength(c)
	case api.APIID:
		return nil, f.handleStateAPIID(c)
	case api.FrameData:
		return nil, f.handleStateFrame(c)
	case api.FrameChecksum:
		return f.handleStateChecksum(c)
	default:
		return nil, f.handleStateStart(c)
	}
}

func (f *APIFrame) handleStateChecksum(c byte) (Frame, error) {
	f.state.state = api.FrameStart
	f.state.chechsum += c

	if api.ValidChecksum != f.state.chechsum {
		return nil, api.ErrChecksumValidation
	}

	return f.frame, nil
}

func (f *APIFrame) handleStateFrame(c byte) error {
	err := f.frame.RX(c)
	if err != nil {
		f.state.state = api.FrameStart
		return err
	}

	f.state.chechsum += c
	f.state.index++

	if f.state.index == f.state.dataSize {
		f.state.state = api.FrameChecksum
	}

	return nil
}

func (f *APIFrame) handleStateAPIID(c byte) error {
	var err error
	f.frame, err = NewFrameForAPIID(c)
	if err != nil {
		f.state.state = api.FrameStart
		return err
	}
	f.state.chechsum += c
	f.state.index++
	f.state.state = api.FrameData

	return nil
}

func (f *APIFrame) handleStateLength(c byte) error {
	f.state.dataSize += uint16(c << (1 - f.state.index))
	f.state.index++

	if f.state.index == api.FrameLengthByteCount {
		f.state.index = 0
		f.state.state = api.APIID
	}
	return nil
}

func (f *APIFrame) handleStateStart(c byte) error {
	if api.FrameDelimiter != c {
		return api.ErrFrameDelimiter
	}

	f.state.escape = false
	f.state.index = 0
	f.state.dataSize = 0
	f.state.chechsum = 0
	f.state.state = api.FrameLength

	return nil
}

func (f *APIFrame) shouldEscapeNext(c byte) bool {
	if f.Mode == api.EscapeModeInactive {
		return false
	}

	if f.state.state == api.FrameStart {
		return false
	}

	if f.state.escape {
		return false
	}

	if c != api.ESC {
		return false
	}

	f.state.escape = true
	return true
}

func (f *APIFrame) escape(c byte) byte {
	f.state.escape = false
	return api.Escape(c)
}

type state struct {
	state    api.State
	escape   bool
	index    uint16
	dataSize uint16
	chechsum uint8
}
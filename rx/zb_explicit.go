package rx

const zbExplicitAPIID byte = 0x91

const (
	zbeAddr64  = rxFrameState(iota)
	zbeAddr16  = rxFrameState(iota)
	zbeSrcEp   = rxFrameState(iota)
	zbeDstEp   = rxFrameState(iota)
	zbeCID     = rxFrameState(iota)
	zbePID     = rxFrameState(iota)
	zbeOptions = rxFrameState(iota)
	zbeData    = rxFrameState(iota)
)

var _ Frame = (*ZBExplicit)(nil)

// ZBExplicit rx frame
type ZBExplicit struct {
	state     rxFrameState
	index     byte
	Addr64    uint64
	Addr16    uint16
	SrcEP     byte
	DstEP     byte
	ClusterID uint16
	ProfileID uint16
	Options   byte
	Data      []byte
}

func newZBExplicit() Frame {
	return &ZBExplicit{
		state: zbeAddr64,
	}
}

// RX frame data
func (f *ZBExplicit) RX(b byte) error {
	var err error

	switch f.state {
	case zbeAddr64:
		err = f.stateAddr64(b)
	case zbeAddr16:
		err = f.stateAddr16(b)
	case zbeSrcEp:
		err = f.stateSrcEP(b)
	case zbeDstEp:
		err = f.stateDstEP(b)
	case zbeCID:
		err = f.stateCID(b)
	case zbePID:
		err = f.statePID(b)
	case zbeOptions:
		err = f.stateOptions(b)
	case zbeData:
		err = f.stateData(b)
	}

	return err
}

func (f *ZBExplicit) stateAddr64(b byte) error {
	f.Addr64 += uint64(b) << (56 - (8 * f.index))
	f.index++

	if f.index == 8 {
		f.index = 0
		f.state = zbeAddr16
	}

	return nil
}

func (f *ZBExplicit) stateAddr16(b byte) error {
	f.Addr16 += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = zbeSrcEp
	}

	return nil

}

func (f *ZBExplicit) stateSrcEP(b byte) error {
	f.SrcEP = b
	f.state = zbeDstEp

	return nil
}

func (f *ZBExplicit) stateDstEP(b byte) error {
	f.DstEP = b
	f.state = zbeCID

	return nil
}

func (f *ZBExplicit) stateCID(b byte) error {
	f.ClusterID += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = zbePID
	}
	return nil
}

func (f *ZBExplicit) statePID(b byte) error {
	f.ProfileID += uint16(b) << (8 - (8 * f.index))
	f.index++

	if f.index == 2 {
		f.index = 0
		f.state = zbeOptions
	}
	return nil
}

func (f *ZBExplicit) stateOptions(b byte) error {
	f.Options = b
	f.state = zbeData

	return nil
}

func (f *ZBExplicit) stateData(b byte) error {
	if f.Data == nil {
		f.Data = make([]byte, 0)
	}

	f.Data = append(f.Data, b)

	return nil
}

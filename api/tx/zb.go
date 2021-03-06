package tx

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const zbAPIID byte = 0x10

// ZB transmit frame
type ZB struct {
	FrameID         byte
	Addr64          uint64
	Addr16          uint16
	BroadcastRadius byte
	Options         byte
	Data            []byte
}

func NewZB(options ...func(interface{})) *ZB {
	f := &ZB{Addr64: 0xFFFF, Addr16: 0xFFFE}

	optionsRunner(f, options...)

	return f
}

// SetFrameID satisfy FrameIDSetter interface
func (f *ZB) SetFrameID(id byte) {
	f.FrameID = id
}

// SetAddr64 satisfy Addr64Setter interface
func (f *ZB) SetAddr64(addr uint64) {
	f.Addr64 = addr
}

// SetAddr16 satisfy Addr16Setter interface
func (f *ZB) SetAddr16(addr uint16) {
	f.Addr16 = addr
}

// SetBroadcastRadius satisfy BroadcastRadiusSetter interface
func (f *ZB) SetBroadcastRadius(hops byte) {
	f.BroadcastRadius = hops
}

// SetOptions satisfy OptionsSetter interface
func (f *ZB) SetOptions(options byte) {
	f.Options = options
}

// SetData satisfy DataSetter interface
func (f *ZB) SetData(data []byte) {
	f.Data = make([]byte, len(data))
	copy(f.Data, data)
}

// Bytes turn ZB frame into bytes, satisfy Frame interface
func (f *ZB) Bytes() ([]byte, error) {
	var b bytes.Buffer

	b.WriteByte(zbAPIID)
	b.WriteByte(f.FrameID)
	b.Write(util.Uint64ToBytes(f.Addr64))
	b.Write(util.Uint16ToBytes(f.Addr16))
	b.WriteByte(f.BroadcastRadius)
	b.WriteByte(f.Options)

	if f.Data != nil && len(f.Data) > 0 {
		b.Write(f.Data)
	}

	return b.Bytes(), nil
}

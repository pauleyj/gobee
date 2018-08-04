package at_remote

import (
	"bytes"
	"github.com/pauleyj/gobee/api/tx/util"
)

const atRemoteAPIID byte = 0x17

func FrameID(id byte) func(*ATRemote) {
	return func(f *ATRemote) {
		f.FrameID = id
	}
}

func Addr64(addr64 uint64) func(*ATRemote) {
	return func(f *ATRemote) {
		f.Addr64 = addr64
	}
}

func Addr16(addr16 uint16) func(*ATRemote) {
	return func(f *ATRemote) {
		f.Addr16 = addr16
	}
}

func Options(options byte) func(*ATRemote) {
	return func(f *ATRemote) {
		f.Options = options
	}
}

func Command(cmd [2]byte) func(*ATRemote) {
	return func(f *ATRemote) {
		copy(f.Cmd[:], cmd[:])
		for i, b := range cmd {
			f.Cmd[i] = b
		}
	}
}

func Parameter(parameter []byte) func(*ATRemote) {
	return func(f *ATRemote) {
		if parameter == nil || len(parameter) == 0 {
			return
		}

		f.Parameter = make([]byte, len(parameter))
		copy(f.Parameter, parameter)
	}
}

func NewATRemote(options ...func(*ATRemote)) *ATRemote {
	f := &ATRemote{Addr64:0xFFFF, Addr16:0xFFFE, Cmd:[2]byte{'N','I'}}

	if options == nil {
		return f
	}

	for _, option := range options {
		option(f)
	}

	return f
}

// ATRemote AT remote transmit frame
type ATRemote struct {
	FrameID byte
	Addr64 uint64
	Addr16 uint16
	Options byte
	Cmd [2]byte
	Parameter []byte
}

// Bytes turn ATRemote frame into bytes
func (f *ATRemote) Bytes() ([]byte, error) {
	var b bytes.Buffer

	err := b.WriteByte(atRemoteAPIID)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.FrameID)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint64ToBytes(f.Addr64))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(util.Uint16ToBytes(f.Addr16))
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(f.Options)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(f.Cmd[:])
	if err != nil {
		return nil, err
	}

	if f.Parameter != nil && len(f.Parameter) > 0 {
		_, err = b.Write(f.Parameter)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

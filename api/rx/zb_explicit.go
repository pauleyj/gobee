package rx

import (
	"encoding/binary"
)

const (
	zbExplicitAPIID byte = 0x91

	zbeAddr64Offset    = 0
	zbeAddr16Offset    = 8
	zbeSrcEPOffset     = 10
	zbeDstEPOffset     = 11
	zbeClusterIDOffset = 12
	zbeClusterIDLength = 2
	zbeProfileIDOffset = 14
	zbeProfileIDLength = 2
	zbeOptionsOffset   = 16
	zbeDataOffset      = 17
)

var _ Frame = (*ZBExplicit)(nil)

// ZBExplicit rx frame
type ZBExplicit struct {
	buffer []byte
}

func newZBExplicit() Frame {
	return &ZBExplicit{
		buffer: make([]byte, 0),
	}
}

// RX frame data
func (f *ZBExplicit) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

// Addr64 64-bit address of sender
func (f *ZBExplicit) Addr64() uint64 {
	return binary.BigEndian.Uint64(f.buffer[zbeAddr64Offset : zbeAddr64Offset+addr64Length])
}

// Addr16 16-bit address of sender
func (f *ZBExplicit) Addr16() uint16 {
	return binary.BigEndian.Uint16(f.buffer[zbeAddr16Offset : zbeAddr16Offset+addr16Length])
}

// SrcEP source endpoint
func (f *ZBExplicit) SrcEP() byte {
	return f.buffer[zbeSrcEPOffset]
}

// DstEP destination endpoint
func (f *ZBExplicit) DstEP() byte {
	return f.buffer[zbeDstEPOffset]
}

// ClusterID clister ID
func (f *ZBExplicit) ClusterID() uint16 {
	return binary.BigEndian.Uint16(f.buffer[zbeClusterIDOffset : zbeClusterIDOffset+zbeClusterIDLength])
}

// ProfileID profile ID
func (f *ZBExplicit) ProfileID() uint16 {
	return binary.BigEndian.Uint16(f.buffer[zbeProfileIDOffset : zbeProfileIDOffset+zbeProfileIDLength])
}

// Options frame options
func (f *ZBExplicit) Options() byte {
	return f.buffer[zbeOptionsOffset]
}

// Data frame data
func (f *ZBExplicit) Data() []byte {
	if len(f.buffer) == zbeDataOffset {
		return nil
	}

	return f.buffer[zbeDataOffset:]
}

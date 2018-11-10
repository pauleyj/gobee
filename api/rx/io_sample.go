package rx

import "encoding/binary"

const (
	ioSampleAPIID byte = 0x92

	ioSampleAddr64Offset            = 0
	ioSampleAddr16Offset            = 8
	ioSampleOptionsOffset           = 9
	ioSampleSampleCountOffset       = 10
	ioSampleDigitalSampleMaskOffset = 11
	ioSampleAnalogSampleMaskOffset  = 13
	ioSampleDigitalSamplesOffset    = 14
	ioSampleAnalogSampleOffset      = 16
)

type IOSample struct {
	buffer []byte
}

func newIOSample() Frame {
	return &IOSample{
		buffer: make([]byte, 0),
	}
}

func (f *IOSample) RX(b byte) error {
	f.buffer = append(f.buffer, b)

	return nil
}

func (f *IOSample) Addr64() uint64 {
	return binary.BigEndian.Uint64(f.buffer[ioSampleAddr64Offset : ioSampleAddr64Offset+addr64Length])
}

func (f *IOSample) Addr16() uint16 {
	return binary.BigEndian.Uint16(f.buffer[ioSampleAddr16Offset : ioSampleAddr16Offset+addr16Length])
}

func (f *IOSample) Options() byte {
	return f.buffer[ioSampleOptionsOffset]
}

func (f *IOSample) SampleCount() byte {
	return f.buffer[ioSampleSampleCountOffset]
}

func (f *IOSample) DigitalSampleMask() uint16 {
	return binary.BigEndian.Uint16(f.buffer[ioSampleDigitalSampleMaskOffset : ioSampleAnalogSampleMaskOffset])
}

func (f *IOSample) AnalogSampleMask() byte {
	return f.buffer[ioSampleAnalogSampleMaskOffset]
}

func (f *IOSample) DigitalSamples() uint16 {
	return binary.BigEndian.Uint16(f.buffer[ioSampleDigitalSamplesOffset:ioSampleAnalogSampleOffset])
}

func (f *IOSample) AnalogSample() uint16 {
	return binary.BigEndian.Uint16(f.buffer[ioSampleAnalogSampleOffset:])
}
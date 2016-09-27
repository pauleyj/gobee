package rx

type rxFrameState int

// Frame interface for RX frames
type Frame interface {
	RX(byte) error
}

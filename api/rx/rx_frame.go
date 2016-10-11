package rx

type rxFrameState int

const (
	addr64Length = 8
	addr16Length = 2
)

// Frame interface for RX frames
type Frame interface {
	RX(byte) error
}

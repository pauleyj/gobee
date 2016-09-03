package rx

type rx_frame_state int

type RxFrame interface {
	RX(byte) error
}

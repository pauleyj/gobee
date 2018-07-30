package tx

// Frame interface for TX frames
type Frame interface {
	Bytes() ([]byte, error)
}

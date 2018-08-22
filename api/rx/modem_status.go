package rx

const (
	modemStatusAPIID byte = 0x8A
)

type ModemStatus struct {
	status byte
}

func newModemStatus() Frame {
	return &ModemStatus{}
}


func (f *ModemStatus) RX(b byte) error {
	f.status = b

	return nil
}

func (f *ModemStatus) Status() byte {
	return f.status
}
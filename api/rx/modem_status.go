package rx

const (
	modemStatusAPIID byte = 0x8A
)

type ModemStatus struct {
	Status byte
}

func newModemStatus() Frame {
	return &ModemStatus{}
}


func (f *ModemStatus) RX(b byte) error {
	f.Status = b

	return nil
}
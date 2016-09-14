package tx

type TxFrame interface {
	Bytes() ([]byte, error)
}

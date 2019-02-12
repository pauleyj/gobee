package common

import (
	"fmt"
	"github.com/tarm/serial"
)

// NewTransmitter constructs a new Transmitter
func NewTransmitter(port *serial.Port, verbose bool) *Transmitter {
	return &Transmitter{
		port:    port,
		verbose: verbose,
	}
}

// Transmitter implements gobee.XBeeTransmitter
type Transmitter struct {
	port    *serial.Port
	verbose bool
}

// Transmit satisifies the gobee.XBeeTransmitter interface.  It transmits
// bytes (tx frame bytes) to the uart.
func (t *Transmitter) Transmit(b []byte) (int, error) {
	if t.verbose {
		fmt.Print("tx --> ")
		for _, b := range b {
			fmt.Printf("%#0.2x ", b)
		}
		fmt.Println()
	}

	return t.port.Write(b)
}

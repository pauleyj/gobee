package client

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/_examples/echo/cmd/common"
	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/rx"
	"github.com/pauleyj/gobee/api/tx"
	"github.com/tarm/serial"
)

func NewEchoClient(port string, baud int, verbose bool) *EchoClient {
	return &EchoClient{
		port:    port,
		baud:    baud,
		verbose: verbose,
		done:    make(chan struct{}),
	}
}

type EchoClient struct {
	port    string
	baud    int
	verbose bool
	done    chan struct{}

	sp *serial.Port
}

func (c *EchoClient) Open() error {
	//
	// configure and open serial port
	cfg := &serial.Config{
		Name:        c.port,
		Baud:        c.baud,
		ReadTimeout: 5 * time.Millisecond,
	}

	var err error
	c.sp, err = serial.OpenPort(cfg)
	if err != nil {
		return err
	}
	//
	// build transmitter
	transmitter := common.NewTransmitter(c.sp, c.verbose)
	//
	// build receiver
	rx := make(chan rx.Frame)
	receiver := common.NewReceiver(rx)
	//
	// build xbee
	xbee := gobee.New(transmitter, receiver, gobee.APIEscapeMode(api.EscapeModeActive))
	//
	// serial port rx loop
	go c.rx(xbee, c.sp)
	//
	// console input loop
	go c.console(xbee)
	//
	// echo client
	go c.client(xbee, rx)

	return nil
}

func (c *EchoClient) Close() error {
	var err error
	//
	// exit client, console, and rx loops
	close(c.done)
	//
	// close serial port
	if c.sp != nil {
		err = c.sp.Close()
	}

	return err
}

func (c *EchoClient) client(xbee *gobee.XBee, ch <-chan rx.Frame) {
	go func() {
		for {
			select {
			case f := <-ch:
				switch f.(type) {
				case *rx.TXStatus:
					fmt.Printf("Echo client received (TXS): %v\n", f.(*rx.TXStatus))
				case *rx.ZB:
					fmt.Printf("Echo client received (ZB): %s\n", string(f.(*rx.ZB).Data()))
				default:
					fmt.Printf("Echo client received unhandled rx frame of type:%+v (%+v)\n", reflect.TypeOf(f), f)
				}
			case <-c.done:
				return
			}
		}
	}()
	//
	// initialize API options and PAN network
	_, err := xbee.TX(tx.NewAT(
		tx.FrameID(1),
		tx.Command([2]byte{'D', 'A'}),
		tx.Parameter([]byte{0x6c, 0x33, 0x33, 0x74})))
	if err != nil {
		fmt.Printf("Failed to set the operating 16-bit PAN ID for the network: %v", err)
	}

	_, err = xbee.TX(tx.NewAT(
		tx.FrameID(2),
		tx.Command([2]byte{'I', 'D'}),
		tx.Parameter([]byte{0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64})))
	if err != nil {
		fmt.Printf("Failed to set the extended PAN ID for the network: %v", err)
	}

	_, err = xbee.TX(tx.NewAT(
		tx.FrameID(3),
		tx.Command([2]byte{'A', 'P'}),
		tx.Parameter([]byte{0x02})))
	if err != nil {
		fmt.Printf("Failed to set the API mode: %v", err)
	}

	fmt.Println("client initialized")

	<-c.done

}

func (c *EchoClient) console(xbee *gobee.XBee) {
	bio := bufio.NewReader(os.Stdin)
	//
	// forever read line from the console, send to coordinator
	for {
		buffer, err := bio.ReadBytes('\n')
		if err != nil {
			fmt.Printf("client failed to read stdio: %v\n", err)
			continue
		}

		if len(buffer) == 1 {
			continue
		}

		payload := buffer[:len(buffer)-1]
		if c.verbose {
			fmt.Printf("client sending data payload: [%s]\n", payload)
		}

		msg := tx.NewZB(
			tx.FrameID(1),
			tx.Addr64(0),
			tx.Addr16(0),
			tx.BroadcastRadius(0),
			tx.Options(0),
			tx.Data(payload))
		_, err = xbee.TX(msg)
		if err != nil {
			fmt.Printf("client failed to transmit frame: %v\n", err)
		}
	}
}

func (c *EchoClient) rx(xbee *gobee.XBee, p *serial.Port) {
	var buf [256]byte
	//
	// forever, RX from uart and process, received API frames
	// will be handled by receiver
	for {
		select {
		case <-c.done:
			return
		default:
			n, _ := p.Read(buf[:])
			if n != 0 {
				if c.verbose {
					fmt.Print("rx <-- ")
				}
				for i := 0; i < n; i++ {
					if c.verbose {
						fmt.Printf("%#0.2x ", buf[i])
					}
					err := xbee.RX(buf[i])
					if err != nil {
						fmt.Printf("\necho failed RX: %v\n", err)
					}
				}
				if c.verbose {
					fmt.Println()
				}
			}
		}
	}
}

package server

import (
	"fmt"
	"github.com/pauleyj/gobee"
	"github.com/pauleyj/gobee/_examples/echo/cmd/common"
	"github.com/pauleyj/gobee/api"
	"github.com/pauleyj/gobee/api/rx"
	"github.com/pauleyj/gobee/api/tx"
	"github.com/tarm/serial"
	"reflect"
	"time"
)

// NewEchoServer constructs a new echo server
func NewEchoServer(port string, baud int, verbose bool) *EchoServer {
	return &EchoServer{
		port:    port,
		baud:    baud,
		verbose: verbose,
		done:    make(chan struct{}),
	}
}

type EchoServer struct {
	port    string
	baud    int
	verbose bool
	done    chan struct{}

	sp   *serial.Port
}

func (s *EchoServer) Open() error {
	//
	// configure and open serial port
	cfg := &serial.Config{
		Name:        s.port,
		Baud:        s.baud,
		ReadTimeout: 1 * time.Millisecond,
	}

	var err error
	s.sp, err = serial.OpenPort(cfg)
	if err != nil {
		return err
	}
	//
	// build transmitter
	transmitter := common.NewTransmitter(s.sp, s.verbose)
	//
	// build receiver
	rx := make(chan rx.Frame)
	receiver := common.NewReceiver(rx)
	//
	// build xbee
	xbee := gobee.New(transmitter, receiver, gobee.APIEscapeMode(api.EscapeModeInactive))
	//
	// serial port rx loop
	go s.rx(xbee, s.sp)
	//
	// serve echo
	go s.serve(xbee, rx)

	return nil
}

func (s *EchoServer) Close() error {
	var err error
	//
	// exit serve and rx loop
	close(s.done)
	//
	// close serial port
	if s.sp != nil {
		err = s.sp.Close()
	}

	return err
}

func (s *EchoServer) serve(xbee *gobee.XBee, ch <-chan rx.Frame) {
	go func() {
		var frameID byte = 0
		for {
			select {
			case f := <-ch:
				switch frame := f.(type) {
				case *rx.ZB:
					fmt.Printf("Echo server received ZB message: '%s' from client: %#0.16x %#0.4x\n", string(frame.Data()), frame.Addr64(), frame.Addr16())
					//
					// echo message back to client
					frameID++
					echo := tx.NewZB(
						tx.FrameID(frameID),
						tx.Addr64(frame.Addr64()),
						tx.Addr16(frame.Addr16()),
						tx.Data(frame.Data()))

					_, err := xbee.TX(echo)
					if err != nil {
						fmt.Printf("Echo server failed to transmit echo: %v\n", err)
					}
				case *rx.TXStatus:
					fmt.Printf("Echo server received TX Status: %v\n", frame)
				case *rx.ModemStatus:
					fmt.Printf("Echo server received modem status: %+v\n", frame)
				default:
					fmt.Printf("Echo server received unhandled rx frame of type:%+v\n", reflect.TypeOf(frame))
				}
				case <-time.After(10 * time.Second):
					fmt.Println("Sending broadcast ping")
					frameID++
					_, err := xbee.TX(tx.NewZB(
						tx.FrameID(frameID),
						tx.Addr64(api.BroadcastAddr64),
						tx.Addr16(api.BroadcastAddr16),
						tx.BroadcastRadius(0),
						tx.Options(0),
						tx.Data([]byte{'p', 'i', 'n', 'g'})))
					if err != nil {
						fmt.Printf("Echo server failed to transmit ping: %v\n", err)
					}
			case <-s.done:
				return
			}
		}
	}()
	//
	// initialize API options and PAN network
	//_, err := xbee.TX(tx.NewAT(
	//	tx.FrameID(1),
	//	tx.Command([2]byte{'F', 'R'}),
	//	tx.Parameter([]byte{0x6c, 0x33, 0x33, 0x74})))
	//if err != nil {
	//	fmt.Printf("Failed to initiate software reset: %v", err)
	//}
	//
	//<- time.After(5 * time.Second)
	//
	//_, err = xbee.TX(tx.NewAT(
	//	tx.FrameID(1),
	//	tx.Command([2]byte{'N', 'R'}),
	//	tx.Parameter([]byte{0x6c, 0x33, 0x33, 0x74})))
	//if err != nil {
	//	fmt.Printf("Failed to initiate network reset: %v", err)
	//}
	//
	//_, err = xbee.TX(tx.NewAT(
	//	tx.FrameID(1),
	//	tx.Command([2]byte{'I', 'I'}),
	//	tx.Parameter([]byte{0x6c, 0x33, 0x33, 0x74})))
	//if err != nil {
	//	fmt.Printf("Failed to set the operating 16-bit PAN ID for the network: %v", err)
	//}

	_, err := xbee.TX(tx.NewAT(
		tx.FrameID(1),
		tx.Command([2]byte{'I', 'D'}),
		tx.Parameter([]byte{0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64})))
	if err != nil {
		fmt.Printf("Failed to set the extended PAN ID for the network: %v", err)
	}

	_, err = xbee.TX(tx.NewAT(
		tx.FrameID(1),
		tx.Command([2]byte{'A', 'P'}),
		tx.Parameter([]byte{0x01})))
	if err != nil {
		fmt.Printf("Failed to set the API mode: %v", err)
	}

	<-s.done
}

func (s *EchoServer) rx(xbee *gobee.XBee, p *serial.Port) {
	var buf [256]byte
	//
	// forever, RX from uart and process, received API frames
	// will be handled by receiver
	for {
		select {
		case <-s.done:
			return
		default:
			n, _ := p.Read(buf[:])
			if n != 0 {
				if s.verbose {
					fmt.Print("rx <-- ")
				}
				for i := 0; i < n; i++ {
					if s.verbose {
						fmt.Printf("%#0.2x ", buf[i])
					}
					err := xbee.RX(buf[i])
					if err != nil {
						fmt.Printf("echo failed RX: %v\n", err)
					}
				}
				if s.verbose {
					fmt.Println()
				}
			}
		}
	}
}
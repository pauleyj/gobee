# gobee - A Go XBee Library

[![Build Status](https://travis-ci.org/pauleyj/gobee.svg?branch=master)](https://travis-ci.org/pauleyj/gobee)
[![Coverage Status](https://coveralls.io/repos/github/pauleyj/gobee/badge.svg)](https://coveralls.io/github/pauleyj/gobee)
[![Go Report Card](https://goreportcard.com/badge/github.com/pauleyj/gobee)](https://goreportcard.com/report/github.com/pauleyj/gobee)
[![codebeat badge](https://codebeat.co/badges/75f31b30-5397-4626-9118-9b599e088f44)](https://codebeat.co/projects/github-com-pauleyj-gobee)

gobee, a library for enabling support of XBee series 2 and series 3 low power radios to your Go project.

---

### Usage

Implement the _XBeeTransmitter_ and _XBeeReceiver_ interfaces, instantiate an XBee, and start communicating, almost... It is up to the gobee user to configure and own the serial port the XBee is connected to and marshal data to/from it to gobee.

#### XBeeTransmitter

```golang
type XBeeTransmitter interface {
	Transmit([]byte) (int, error)
}
```

gobee uses the XBeeTransmitter interface to request a byte slice be sent to the serial communications port the XBee is connected to.

#### XBeeReceiver

```golang
type XBeeReceiver interface {
	Receive(rx.RxFrame) error
}
```

gobee uses the XBeeReceiver interface to report received API frames.

#### XBee

This the XBee widget used to communicate with the physical XBee.

```golang
...
transmitter := &Transmitter{...}	// your XBeeTransmitter
receiver    := &Receiver{...}		// your XBeeReceiver
xbee        := gobee.New(transmitter, receiver)
...
```

#### Transmitting an API Frame

To send an API frame, construct the frame using an appropriate constructor and option functions to set frame parameters.  

```golang
frame  := tx.NewZB(
			tx.FrameID(frameID),
			tx.Addr64(api.BroadcastAddr64),
			tx.Addr16(api.BroadcastAddr16),
			tx.Data(frame.Data()))
```

gobee will then call the supplied Transmitters Transmit to send the bytes to the UART the XBee is connected to.

```golang
at := tx.NewATBuilder().
		ID(1).
		Command([2]byte{'A','O'}).
		Parameter(nil).
		Build()
_, err := xbee.TX(at)
...
```

#### Sending API Frame to the UART

When a frame is transmitted, gobee forms an appropriate API packet (see Transmitting a Frame) and sends it to your XBeeTransmitter for writing to the serial UART the XBee is connected to.

```golang
...
func (tx *Transmitter) Transmit(buffer []byte) (n int, err error) {

	i, err := port.Write(buffer)
	if err != nil {
		fmt.Printf("Failed to write buffer to xbee comms: %v\n", err)
	}

	return i, err
}
...
```

#### Receiving Data from the UART

When data is received from the serial UART the XBee is connected to, send it to gobee.

```golang
...
n, err := port.Read(buffer)
...
for i := 0; i < n; i++ {
	err = xbee.RX(buffer[i])
	...
}
...
```

#### Receiving Data Frames from gobee

Your implemented XBeeReceiver will get called when completed API frames are received.  gobee validates the received API frame and reports the data frame via the XBeeReceiver interface.

```golang
func (r *Receiver) Receive(f rx.Frame) error {
	switch f.(type) {
	case *rx.ZB:
		// do something with received ZB frame
	...
	}

	return nil
}
```

### License

gobee is licensed under the MIT License.  See the [LICENSE](https://github.com/pauleyj/gobee/blob/master/LICENSE) for more information.
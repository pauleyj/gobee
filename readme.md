# gobee - A Go XBee Library

[![Build Status](https://travis-ci.org/pauleyj/gobee.svg?branch=master)](https://travis-ci.org/pauleyj/gobee)

gobee, a library for enabling support of XBee series 2 and series 3 low power radios to your Go project.

---

### Usage

Implement the _XBeeTransmitter_ and _XBeeReceiver_ interfaces, instantiate an XBee, and start communicating, almost... It is up to the gobee user to configure and own the serial port the XBee is connected to and marshal data to/from it to gobee.

#### XBeeTransmitter

```golang
type XBeeTransmitter interface {
	io.Writer
}
```

gobee uses the XBeeTransmitter interface to request a byte slice be sent to the serial communications port the XBee is connected to.

#### XBeeReceiver

```golang
type XBeeReceiver interface {
	OnRxFrame(rx.RxFrame) error
}
```

gobee uses the XBeeReceiver interface to report received API frames.

#### XBee

This the XBee widget used to communicate with the physical XBee.

```golang
...
transmitter := &Transmitter{...}	// your XBeeTransmitter
receiver    := &Receiver{...}		// your XBeeReceiver
xbee        := gobee.NewXBee(transmitter, receiver)
...
```

#### Sending Data to the UART

When a frame is transmitted, gobee forms an appropriate API packet and sends it to your XBeeTransmitter for writing to the serial UART the physical XBee is connected to.

```golang
...
func (tx *Transmitter) Write(buffer []byte) (n int, err error) {

	i, err := port.Write(buffer)
	if err != nil {
		fmt.Printf("Failed to write buffer to xbee comms: %v\n", err)
	}

	return i, err
}
...
```

#### Receiving Data from the UART

When data is received from the serial UART the physical XBee is connected to, send it to gobee.

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

#### Receiving API Frames from gobee

Your implemented XBeeReceiver will get called when completed API frames are received.

```golang
func (r *Receiver) OnRxFrame(f rx.RxFrame) error {
	switch f.(type) {
	case *rx.ZB:
		// do something with received ZB frame
	...
	}

	return nil
}
```
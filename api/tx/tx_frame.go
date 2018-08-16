package tx

// Frame interface for TX frames
type Frame interface {
	Bytes() ([]byte, error)
}

type FrameIDSetter interface {
	SetFrameID(byte)
}

func FrameID(id byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(FrameIDSetter); ok {
			f.SetFrameID(id)
		}
	}
}

type CommandSetter interface {
	SetCommand([2]byte)
}

func Command(cmd [2]byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(CommandSetter); ok {
			f.SetCommand(cmd)
		}
	}
}

type ParameterSetter interface {
	SetParameter([]byte)
}

func Parameter(parameter []byte) func(interface{}) {
	return func(i interface{}) {
		if parameter == nil || len(parameter) == 0 {
			return
		}

		if f, ok := i.(ParameterSetter); ok {
			f.SetParameter(parameter)
		}
	}
}

type Addr64Setter interface {
	SetAddr64(uint64)
}

func Addr64(addr uint64) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(Addr64Setter); ok {
			f.SetAddr64(addr)
		}
	}
}

type Addr16Setter interface {
	SetAddr16(uint16)
}

func Addr16(addr uint16) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(Addr16Setter); ok {
			f.SetAddr16(addr)
		}
	}
}

type OptionsSetter interface {
	SetOptions(byte)
}

func Options(options byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(OptionsSetter); ok {
			f.SetOptions(options)
		}
	}
}

type BroadcastRadiusSetter interface {
	SetBroadcastRadius(byte)
}

func BroadcastRadius(hops byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(BroadcastRadiusSetter); ok {
			f.SetBroadcastRadius(hops)
		}
	}
}

type DataSetter interface {
	SetData([]byte)
}

func Data(data []byte) func(interface{}) {
	return func(i interface{}) {
		if data == nil || len(data) == 0 {
			return
		}

		if f, ok := i.(DataSetter); ok {
			f.SetData(data)
		}
	}
}

type SrcEPSetter interface {
	SetSrcEP(byte)
}

func SrcEP(src byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(SrcEPSetter); ok {
			f.SetSrcEP(src)
		}
	}
}

type DstEPSetter interface {
	SetDstEP(byte)
}

func DstEP(dst byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(DstEPSetter); ok {
			f.SetDstEP(dst)
		}
	}
}

type ClusterIDSetter interface {
	SetClusterID(byte)
}

func ClusterID(id byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(ClusterIDSetter); ok {
			f.SetClusterID(id)
		}
	}
}

type ProfileIDSetter interface {
	SetProfileID(byte)
}

func ProfileID(id byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(ProfileIDSetter); ok {
			f.SetProfileID(id)
		}
	}
}

func optionsRunner(i interface{}, options ...func(interface{})) {
	if options == nil || len(options) == 0 {
		return
	}

	for _, option := range options {
		if option == nil {
			continue
		}

		option(i)
	}
}
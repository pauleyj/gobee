package tx

// Frame interface for TX frames
type Frame interface {
	Bytes() ([]byte, error)
}

// FrameIDSetter sets the frame ID
type FrameIDSetter interface {
	SetFrameID(byte)
}

// FrameID helper options function to set frame ID
func FrameID(id byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(FrameIDSetter); ok {
			f.SetFrameID(id)
		}
	}
}

// CommandSetter sets AT related frame command
type CommandSetter interface {
	SetCommand([2]byte)
}

// Command helper options function to set AT related frame command
func Command(cmd [2]byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(CommandSetter); ok {
			f.SetCommand(cmd)
		}
	}
}

// ParameterSetter sets frame parameters
type ParameterSetter interface {
	SetParameter([]byte)
}

// Parameter helper options function to set frame parameter
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

// Addr64Setter sets frame 64-bit address
type Addr64Setter interface {
	SetAddr64(uint64)
}

// Addr64 helper options function to set frames 64-bit address
func Addr64(addr uint64) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(Addr64Setter); ok {
			f.SetAddr64(addr)
		}
	}
}

// Addr64Setter sets frame 16-bit address
type Addr16Setter interface {
	SetAddr16(uint16)
}

// Addr16 helper options function to set frames 16-bit address
func Addr16(addr uint16) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(Addr16Setter); ok {
			f.SetAddr16(addr)
		}
	}
}

// OptionsSetter sets frame options
type OptionsSetter interface {
	SetOptions(byte)
}

// Options helper options function to set frame options
func Options(options byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(OptionsSetter); ok {
			f.SetOptions(options)
		}
	}
}

// BroadcastRadiusSetter sets frame broadcast radius
type BroadcastRadiusSetter interface {
	SetBroadcastRadius(byte)
}

// BroadcastRadius helper options function to set frame broadcast radius
func BroadcastRadius(hops byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(BroadcastRadiusSetter); ok {
			f.SetBroadcastRadius(hops)
		}
	}
}

// DataSetter sets frame data
type DataSetter interface {
	SetData([]byte)
}

// Data helper options function to set frame data
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

// SrcEPSetter sets frame source endpoint
type SrcEPSetter interface {
	SetSrcEP(byte)
}

// SrcEP helper options function to set frame source endpoint
func SrcEP(src byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(SrcEPSetter); ok {
			f.SetSrcEP(src)
		}
	}
}

// DstEPSetter sets frame destination endpoint
type DstEPSetter interface {
	SetDstEP(byte)
}

// DstEP helper options function to set frame destination endpoint
func DstEP(dst byte) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(DstEPSetter); ok {
			f.SetDstEP(dst)
		}
	}
}

// ClusterIDSetter sets frame cluster ID
type ClusterIDSetter interface {
	SetClusterID(uint16)
}

// ClusterID helper options function to set frame cluster ID
func ClusterID(id uint16) func(interface{}) {
	return func(i interface{}) {
		if f, ok := i.(ClusterIDSetter); ok {
			f.SetClusterID(id)
		}
	}
}

// ProfileIDSetter sets frame profile ID
type ProfileIDSetter interface {
	SetProfileID(uint16)
}

// ProfileID helper options function to set frame profile ID
func ProfileID(id uint16) func(interface{}) {
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

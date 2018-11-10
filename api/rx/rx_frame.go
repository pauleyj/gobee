package rx

type rxFrameState int

const (
	addr64Length = 8
	addr16Length = 2
)

// Frame interface for RX frames
type Frame interface {
	RX(byte) error
}

// IDGetter gets frame ID
type IDGetter interface {
	ID() byte
}

// Addr64Getter gets 64-address
type Addr64Getter interface {
	Addr64() uint64
}

// Addr16Getter get 16-bit address
type Addr16Getter interface {
	Addr16() uint16
}

// SrcEPGetter gets source endpoint
type SrcEPGetter interface {
	SrcEP() byte
}
// DstEPGetter getsdestination endpoint
type DstEPGetter interface {
	DstEP() byte
}

// ClusterIDGetter gets cluster ID
type ClusterIDGetter interface {
	ClusterID() uint16
}

// ProfileIDGetter gets profile ID
type ProfileIDGetter interface {
	ProfileID() uint16
}

// CommandGetter gets command
type CommandGetter interface {
	Command() []byte
}

// StatusGetter gets status
type StatusGetter interface {
	Status() byte
}

// OptionsGetter gets options
type OptionsGetter interface {
	Options() byte
}

// DataGetter gets data
type DataGetter interface {
	Data() []byte
}

// RetriesGetter gets retries
type RetriesGetter interface {
	Retries() byte
}

// DeliveryGetter gets delivery
type DeliveryGetter interface {
	Delivery() byte
}

// DiscoveryGetter gets discovery
type DiscoveryGetter interface {
	Discovery() byte
}

// SampleCountGetter gets sample count
type SampleCountGetter interface {
	SampleCount() byte
}

// DigitalSampleMaskGetter gets digital sample mask
type DigitalSampleMaskGetter interface {
	DigitalSampleMask() uint16
}

// AnalogSampleMaskGetter gets analog sample mask
type AnalogSampleMaskGetter interface {
	AnalogSampleMask() byte
}

// DigitalSamplesGetter gets digital samples
type DigitalSamplesGetter interface {
	DigitalSamples() uint16
}

// AnalogSampleGetter gets analog sample
type AnalogSampleGetter interface {
	AnalogSample() uint16
}
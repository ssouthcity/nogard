package nogard

type Availability uint8

const (
	Permanent Availability = iota
	Available
	Unavailable
)

//go:generate stringer -type Availability

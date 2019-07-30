package erb

// FixType represents a type of navigation fix.
type FixType uint8

const (
	FixTypeNoFix  FixType = 0x00
	FixTypeSingle FixType = 0x01
	FixTypeFloat  FixType = 0x02
	FixTypeRTK    FixType = 0x03
)

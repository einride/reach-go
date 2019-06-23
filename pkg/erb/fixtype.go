package erb

// FixType represents a type of navigation fix.
type FixType uint8

//go:generate gobin -m -run golang.org/x/tools/cmd/stringer -type=FixType -trimprefix=FixType

const (
	FixTypeNoFix  FixType = 0x00
	FixTypeSingle FixType = 0x01
	FixTypeFloat  FixType = 0x02
	FixTypeRTK    FixType = 0x03
)

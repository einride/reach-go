package erb

// SVType represents the type of an SV (space vehicle).
type SVType uint8

//go:generate gobin -m -run golang.org/x/tools/cmd/stringer -type=SVType -trimprefix=SVType

const (
	SVTypeGPS     SVType = 0
	SVTypeGLONASS SVType = 1
	SVTypeGalileo SVType = 2
	SVTypeQZSS    SVType = 3
	SVTypeBeiDou  SVType = 4
	SVTypeLEO     SVType = 5
	SVTypeSBAS    SVType = 6
)

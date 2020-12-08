package erb

import (
	"encoding/binary"
	"fmt"
)

// structure of DOPS message.
const (
	lengthOfDOPS           = 12
	indexOfDOPSGeo         = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfDOPSGeo        = 2
	indexOfDOPSPosition    = indexOfDOPSGeo + lengthOfDOPSGeo
	lengthOfDOPSPosition   = 2
	indexOfDOPSVertical    = indexOfDOPSPosition + lengthOfDOPSPosition
	lengthOfDOPSVertical   = 2
	indexOfDOPSHorizontal  = indexOfDOPSVertical + lengthOfDOPSVertical
	lengthOfDOPSHorizontal = 2
)

// compile-time assertion on structure of DOPS message.
var _ [lengthOfDOPS]struct{} = [indexOfDOPSHorizontal + lengthOfDOPSHorizontal]struct{}{}

// scaleOfDOPS is the scale factor to multiply raw DOP values with.
const scaleOfDOPS = 0.01

// DOPS message outputs dimensionless values of DOP.
//
// The raw values are scaled by factor 100.
//
// For example, if received value is 123, then real is 1.23.
type DOPS struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// Geometric DOP.
	Geometric float64
	// Position DOP.
	Position float64
	// Vertical DOP.
	Vertical float64
	// Horizontal DOP.
	Horizontal float64
}

func (d *DOPS) unmarshal(b []byte) error {
	if len(b) != lengthOfDOPS {
		return fmt.Errorf("unmarshal DOPS: unexpected length: %d, expected: %d", len(b), lengthOfDOPS)
	}
	d.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	d.Geometric = scaleOfDOPS * float64(
		binary.LittleEndian.Uint16(b[indexOfDOPSGeo:indexOfDOPSGeo+lengthOfDOPSGeo]),
	)
	d.Position = scaleOfDOPS * float64(
		binary.LittleEndian.Uint16(b[indexOfDOPSPosition:indexOfDOPSPosition+lengthOfDOPSPosition]),
	)
	d.Vertical = scaleOfDOPS * float64(
		binary.LittleEndian.Uint16(b[indexOfDOPSVertical:indexOfDOPSVertical+lengthOfDOPSVertical]),
	)
	d.Horizontal = scaleOfDOPS * float64(
		binary.LittleEndian.Uint16(b[indexOfDOPSHorizontal:indexOfDOPSHorizontal+lengthOfDOPSHorizontal]),
	)
	return nil
}

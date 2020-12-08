package erb

import (
	"encoding/binary"
	"fmt"
)

// structure of VEL message.
const (
	lengthOfVEL              = 28
	indexOfVELNorth          = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfVELNorth         = 4
	indexOfVELEast           = indexOfVELNorth + lengthOfVELNorth
	lengthOfVELEast          = 4
	indexOfVELDown           = indexOfVELEast + lengthOfVELEast
	lengthOfVELDown          = 4
	indexOfVELSpeed          = indexOfVELDown + lengthOfVELDown
	lengthOfVELSpeed         = 4
	indexOfVELHeading        = indexOfVELSpeed + lengthOfVELSpeed
	lengthOfVELHeading       = 4
	indexOfVELSpeedAccuracy  = indexOfVELHeading + lengthOfVELHeading
	lengthOfVELSpeedAccuracy = 4
)

// scaling constants of VEL message.
const scalingOfVELHeading = 1e-5

// compile-time assertion on structure of VEL message.
var _ [lengthOfVEL]struct{} = [indexOfVELSpeedAccuracy + lengthOfVELSpeedAccuracy]struct{}{}

// VEL message contains the velocity in NED (North East Down) coordinates.
type VEL struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// North velocity component (cm/s).
	NorthCentimetersPerSecond int32
	// East velocity component (cm/s).
	EastCentimetersPerSecond int32
	// Down velocity component (cm/s).
	DownCentimetersPerSecond int32
	// Speed is the 2D ground speed (cm/s).
	SpeedCentimetersPerSecond int32
	// Heading is the 2D heading of motion.
	HeadingDegrees float64
	// SpeedAccuracy is the speed accuracy estimate.
	SpeedAccuracyCentimetersPerSecond uint32
}

func (v *VEL) unmarshal(b []byte) error {
	if len(b) != lengthOfVEL {
		return fmt.Errorf("unmarshal VEL: unexpected length: %d, expected: %d", len(b), lengthOfVEL)
	}
	v.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	v.NorthCentimetersPerSecond = int32(binary.LittleEndian.Uint32(
		b[indexOfVELNorth : indexOfVELNorth+lengthOfVELNorth],
	))
	v.EastCentimetersPerSecond = int32(binary.LittleEndian.Uint32(
		b[indexOfVELEast : indexOfVELEast+lengthOfVELEast],
	))
	v.DownCentimetersPerSecond = int32(binary.LittleEndian.Uint32(
		b[indexOfVELDown : indexOfVELDown+lengthOfVELDown],
	))
	v.SpeedCentimetersPerSecond = int32(binary.LittleEndian.Uint32(
		b[indexOfVELSpeed : indexOfVELSpeed+lengthOfVELSpeed],
	))
	v.HeadingDegrees = float64(
		int32(binary.LittleEndian.Uint32(
			b[indexOfVELHeading:indexOfVELHeading+lengthOfVELHeading],
		)),
	) * scalingOfVELHeading
	v.SpeedAccuracyCentimetersPerSecond = binary.LittleEndian.Uint32(
		b[indexOfVELSpeedAccuracy : indexOfVELSpeedAccuracy+lengthOfVELSpeedAccuracy],
	)
	return nil
}

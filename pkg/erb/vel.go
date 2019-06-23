package erb

import (
	"encoding/binary"

	"github.com/einride/unit"
	"golang.org/x/xerrors"
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
	// North velocity component.
	North unit.Speed
	// East velocity component.
	East unit.Speed
	// Down velocity component.
	Down unit.Speed
	// Speed is the 2D ground speed.
	Speed unit.Speed
	// Heading is the 2D heading of motion.
	Heading unit.Angle
	// SpeedAccuracy is the speed accuracy estimate.
	SpeedAccuracy unit.Speed
}

func (v *VEL) unmarshal(b []byte) error {
	if len(b) != lengthOfVEL {
		return xerrors.Errorf("unmarshal VEL: unexpected length: %d, expected: %d", len(b), lengthOfVEL)
	}
	v.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	v.North = unit.Speed(
		int32(binary.LittleEndian.Uint32(
			b[indexOfVELNorth:indexOfVELNorth+lengthOfVELNorth],
		)),
	) * unit.Centi * unit.MetrePerSecond
	v.East = unit.Speed(
		int32(binary.LittleEndian.Uint32(
			b[indexOfVELEast:indexOfVELEast+lengthOfVELEast],
		)),
	) * unit.Centi * unit.MetrePerSecond
	v.Down = unit.Speed(
		int32(binary.LittleEndian.Uint32(
			b[indexOfVELDown:indexOfVELDown+lengthOfVELDown],
		)),
	) * unit.Centi * unit.MetrePerSecond
	v.Speed = unit.Speed(
		binary.LittleEndian.Uint32(
			b[indexOfVELSpeed:indexOfVELSpeed+lengthOfVELSpeed],
		),
	) * unit.Centi * unit.MetrePerSecond
	v.Heading = unit.Angle(
		int32(binary.LittleEndian.Uint32(
			b[indexOfVELHeading:indexOfVELHeading+lengthOfVELHeading],
		)),
	) * unit.Degree * scalingOfVELHeading
	v.SpeedAccuracy = unit.Speed(
		binary.LittleEndian.Uint32(
			b[indexOfVELSpeedAccuracy:indexOfVELSpeedAccuracy+lengthOfVELSpeedAccuracy],
		),
	) * unit.Centi * unit.MetrePerSecond
	return nil
}

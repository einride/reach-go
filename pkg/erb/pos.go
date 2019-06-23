package erb

import (
	"encoding/binary"
	"math"

	"github.com/einride/unit"
	"golang.org/x/xerrors"
)

// structure of POS message.
const (
	lengthOfPOS                     = 44
	indexOfPOSLongitude             = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfPOSLongitude            = 8
	indexOfPOSLatitude              = indexOfPOSLongitude + lengthOfPOSLongitude
	lengthOfPOSLatitude             = 8
	indexOfPOSAltitudeEllipsoid     = indexOfPOSLatitude + lengthOfPOSLatitude
	lengthOfPOSAltitudeEllipsoid    = 8
	indexOfPOSAltitudeMeanSeaLevel  = indexOfPOSAltitudeEllipsoid + lengthOfPOSAltitudeEllipsoid
	lengthOfPOSAltitudeMeanSeaLevel = 8
	indexOfPOSHorizontalAccuracy    = indexOfPOSAltitudeMeanSeaLevel + lengthOfPOSAltitudeMeanSeaLevel
	lengthOfPOSHorizontalAccuracy   = 4
	indexOfPOSVerticalAccuracy      = indexOfPOSHorizontalAccuracy + lengthOfPOSHorizontalAccuracy
	lengthOfPOSVerticalAccuracy     = 4
)

// compile-time assertion on structure of POS packet
var _ [lengthOfPOS]struct{} = [indexOfPOSVerticalAccuracy + lengthOfPOSVerticalAccuracy]struct{}{}

// POS message contains the geodetic coordinates.
//
// Longitude, latitude, altitude and information about accuracy estimate.
type POS struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// Longitude component.
	Longitude float64
	// Latitude component.
	Latitude float64
	// AltitudeEllipsoid is the height above ellipsoid.
	AltitudeEllipsoid unit.Distance
	// AltitudeMeanSeaLevel is the height above mean sea level.
	AltitudeMeanSeaLevel unit.Distance
	// HorizontalAccuracy is the horizontal accuracy estimate in millimeters.
	HorizontalAccuracy unit.Distance
	// VerticalAccuracy is the vertical accuracy estimate in millimeters.
	VerticalAccuracy unit.Distance
}

func (p *POS) unmarshal(b []byte) error {
	if len(b) != lengthOfPOS {
		return xerrors.Errorf("unmarshal POS: unexpected length: %d, expected: %d", len(b), lengthOfPOS)
	}
	p.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	p.Longitude = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSLongitude : indexOfPOSLongitude+lengthOfPOSLongitude],
		),
	)
	p.Latitude = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSLatitude : indexOfPOSLatitude+lengthOfPOSLatitude],
		),
	)
	p.AltitudeEllipsoid = unit.Distance(
		math.Float64frombits(
			binary.LittleEndian.Uint64(
				b[indexOfPOSAltitudeEllipsoid:indexOfPOSAltitudeEllipsoid+lengthOfPOSAltitudeEllipsoid],
			),
		),
	) * unit.Metre
	p.AltitudeMeanSeaLevel = unit.Distance(
		math.Float64frombits(
			binary.LittleEndian.Uint64(
				b[indexOfPOSAltitudeMeanSeaLevel:indexOfPOSAltitudeMeanSeaLevel+lengthOfPOSAltitudeMeanSeaLevel],
			),
		),
	) * unit.Metre
	p.HorizontalAccuracy = unit.Distance(
		binary.LittleEndian.Uint32(
			b[indexOfPOSHorizontalAccuracy:indexOfPOSHorizontalAccuracy+lengthOfPOSHorizontalAccuracy],
		),
	) * unit.Milli * unit.Metre
	p.VerticalAccuracy = unit.Distance(
		binary.LittleEndian.Uint32(
			b[indexOfPOSVerticalAccuracy:indexOfPOSVerticalAccuracy+lengthOfPOSVerticalAccuracy],
		),
	) * unit.Milli * unit.Metre
	return nil
}

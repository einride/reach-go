package erb

import (
	"encoding/binary"
	"math"
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
	// Longitude component (degrees).
	LongitudeDegrees float64
	// Latitude component (degrees).
	LatitudeDegrees float64
	// AltitudeEllipsoid is the height above ellipsoid (m).
	AltitudeEllipsoidMeters float64
	// AltitudeMeanSeaLevel is the height above mean sea level (m).
	AltitudeMeanSeaLevelMeters float64
	// HorizontalAccuracyMillimeters is the horizontal accuracy estimate (mm).
	HorizontalAccuracyMillimeters uint32
	// VerticalAccuracy is the vertical accuracy estimate (mm).
	VerticalAccuracyMillimeters uint32
}

func (p *POS) unmarshalPayload(b []byte) {
	_ = b[lengthOfPOS-1] // early bounds check
	p.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	p.LongitudeDegrees = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSLongitude : indexOfPOSLongitude+lengthOfPOSLongitude],
		),
	)
	p.LatitudeDegrees = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSLatitude : indexOfPOSLatitude+lengthOfPOSLatitude],
		),
	)
	p.AltitudeEllipsoidMeters = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSAltitudeEllipsoid : indexOfPOSAltitudeEllipsoid+lengthOfPOSAltitudeEllipsoid],
		),
	)
	p.AltitudeMeanSeaLevelMeters = math.Float64frombits(
		binary.LittleEndian.Uint64(
			b[indexOfPOSAltitudeMeanSeaLevel : indexOfPOSAltitudeMeanSeaLevel+lengthOfPOSAltitudeMeanSeaLevel],
		),
	)
	p.HorizontalAccuracyMillimeters = binary.LittleEndian.Uint32(
		b[indexOfPOSHorizontalAccuracy : indexOfPOSHorizontalAccuracy+lengthOfPOSHorizontalAccuracy],
	)
	p.VerticalAccuracyMillimeters = binary.LittleEndian.Uint32(
		b[indexOfPOSVerticalAccuracy : indexOfPOSVerticalAccuracy+lengthOfPOSVerticalAccuracy],
	)
}

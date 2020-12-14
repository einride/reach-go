package erb

import (
	"encoding/binary"
)

// structure of STAT message.
const (
	lengthOfSTAT              = 9
	indexOfSTATWeekGPS        = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfSTATWeekGPS       = 2
	indexOfSTATFixType        = indexOfSTATWeekGPS + lengthOfSTATWeekGPS
	lengthOfSTATFixType       = 1
	indexOfSTATHasFix         = indexOfSTATFixType + lengthOfSTATFixType
	lengthOfSTATHasFix        = 1
	indexOfSTATNumSatellites  = indexOfSTATHasFix + lengthOfSTATHasFix
	lengthOfSTATNumSatellites = 1
)

// compile-time assertion on structure of STAT message.
var _ [lengthOfSTAT]struct{} = [indexOfSTATNumSatellites + lengthOfSTATNumSatellites]struct{}{}

// STAT message contains status of fix, its type and also the number of used satellites.
type STAT struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// WeekGPS is the week number of the navigation epoch.
	WeekGPS uint16
	// FixType is the fix type.
	FixType FixType
	// HasFix is true when position and velocity are valid.
	HasFix bool
	// NumSVs is the number of used space vehicles.
	NumSVs uint8
}

func (s *STAT) unmarshalPayload(b []byte) {
	s.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	s.WeekGPS = binary.LittleEndian.Uint16(b[indexOfSTATWeekGPS : indexOfSTATWeekGPS+lengthOfSTATWeekGPS])
	s.FixType = FixType(b[indexOfSTATFixType])
	s.HasFix = b[indexOfSTATHasFix] == 1
	s.NumSVs = b[indexOfSTATNumSatellites]
}

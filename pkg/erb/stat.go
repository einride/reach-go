package erb

import (
	"encoding/binary"
	"fmt"
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

func (s *STAT) unmarshal(b []byte) error {
	if len(b) != lengthOfSTAT {
		return fmt.Errorf("unmarshal STAT: unexpected length: %d, expected: %d", len(b), lengthOfSTAT)
	}
	s.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	s.WeekGPS = binary.LittleEndian.Uint16(b[indexOfSTATWeekGPS : indexOfSTATWeekGPS+lengthOfSTATWeekGPS])
	s.FixType = FixType(b[indexOfSTATFixType])
	switch b[indexOfSTATHasFix] {
	case 0x00:
		s.HasFix = false
	case 0x01:
		s.HasFix = true
	default:
		return fmt.Errorf("unmarshal status: unexpected value of fix status: %d", b[indexOfSTATHasFix])
	}
	s.NumSVs = b[indexOfSTATNumSatellites]
	return nil
}

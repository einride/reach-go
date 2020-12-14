package erb

import (
	"encoding/binary"
)

// structure of VER message.
const (
	lengthOfVER       = 7
	indexOfVERHigh    = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfVERHigh   = 1
	indexOfVERMedium  = indexOfVERHigh + lengthOfVERHigh
	lengthOfVERMedium = 1
	indexOfVERLow     = indexOfVERMedium + lengthOfVERMedium
	lengthOfVERLow    = 1
)

// compile-time assertion on structure of VER message.
var _ [lengthOfVER]struct{} = [indexOfVERLow + lengthOfVERLow]struct{}{}

// VER message contains version of the ERB protocol.
//
// It comprises 3 numbers: high level of version, medium level of version and low level of version.
type VER struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// High level of version.
	High uint8
	// Medium level of version.
	Medium uint8
	// Low level of version.
	Low uint8
}

func (v *VER) unmarshalPayload(b []byte) {
	_ = b[lengthOfVER-1] // early bounds check
	v.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	v.High = b[indexOfVERHigh]
	v.Medium = b[indexOfVERMedium]
	v.Low = b[indexOfVERLow]
}

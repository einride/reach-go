package erb

import (
	"encoding/binary"
	"fmt"

	"github.com/einride/unit"
)

// structure of SVI message.
const (
	lengthOfSVI                   = 25
	lengthOfSV                    = 20
	indexOfNumSVs                 = indexOfTimeGPS + lengthOfTimeGPS
	lengthOfNumSVs                = 1
	indexOfSV                     = indexOfNumSVs + lengthOfNumSVs
	indexOfSVID                   = indexOfSV
	lengthOfSVID                  = 1
	indexOfSVType                 = indexOfSVID + lengthOfSVID
	lengthOfSVType                = 1
	indexOfSVCarrierPhase         = indexOfSVType + lengthOfSVType
	lengthOfSVCarrierPhase        = 4
	indexOfSVPseudoRangeResidual  = indexOfSVCarrierPhase + lengthOfSVCarrierPhase
	lengthOfSVPseudoRangeResidual = 4
	indexOfSVDopplerFrequency     = indexOfSVPseudoRangeResidual + lengthOfSVPseudoRangeResidual
	lengthOfSVDopplerFrequency    = 4
	indexOfSVSignalStrength       = indexOfSVDopplerFrequency + lengthOfSVDopplerFrequency
	lengthOfSVSignalStrength      = 2
	indexOfSVAzimuth              = indexOfSVSignalStrength + lengthOfSVSignalStrength
	lengthOfSVAzimuth             = 2
	indexOfSVElevation            = indexOfSVAzimuth + lengthOfSVAzimuth
	lengthOfSVElevation           = 2
)

// scaling constants of SVI message.
const (
	scaleOfSVCarrierPhase     = 1e-2
	scaleOfSVDopplerFrequency = 1e-3
	scaleOfSVSignalStrength   = 0.25
	scaleOfSVAzimuth          = 1e-1
	scaleOfSVElevation        = 1e-1
)

// compile-time assertions on structure of SVI message.
var (
	_ [lengthOfSVI]struct{} = [indexOfSVElevation + lengthOfSVElevation]struct{}{}
	_ [lengthOfSVI]struct{} = [indexOfNumSVs + lengthOfNumSVs + lengthOfSV]struct{}{}
)

// SVI message contains information about used observation satellites.
type SVI struct {
	// TimeGPS is the time of week in milliseconds of the navigation epoch.
	TimeGPS uint32
	// NumSVs is the number of visible SVs.
	NumSVs uint8
}

func (s *SVI) unmarshal(b []byte) error {
	const expectedLength = indexOfNumSVs + lengthOfNumSVs
	if len(b) < expectedLength {
		return fmt.Errorf("unmarshal SVI: unexpected length: %d, expected: %d", len(b), expectedLength)
	}
	s.TimeGPS = binary.LittleEndian.Uint32(b[indexOfTimeGPS : indexOfTimeGPS+lengthOfTimeGPS])
	s.NumSVs = b[indexOfNumSVs]
	return nil
}

// SV message contains information about a single observation satellite.
type SV struct {
	// ID of SV.
	ID uint8
	// Type of SV.
	Type SVType
	// SignalStrength of SV in dB-Hz.
	SignalStrength float64
	// CarrierPhase of SV in cycles.
	CarrierPhase float64
	// PseudoRangeResidual of SV.
	PseudoRangeResidual unit.Distance
	// DopplerFrequencyHz of SV.
	DopplerFrequencyHz float64
	// Azimuth of SV.
	Azimuth unit.Angle
	// Elevation of SV.
	Elevation unit.Angle
}

func (s *SV) unmarshal(b []byte, i int) error {
	offset := i * lengthOfSV
	expectedLength := indexOfSV + offset + lengthOfSV
	if len(b) < expectedLength {
		return fmt.Errorf("unmarshal SV: unexpected length: %d, expected: %d", len(b), expectedLength)
	}
	s.ID = b[offset+indexOfSVID]
	s.Type = SVType(b[offset+indexOfSVType])
	s.SignalStrength = scaleOfSVSignalStrength * float64(binary.LittleEndian.Uint16(
		b[offset+indexOfSVSignalStrength:offset+indexOfSVSignalStrength+lengthOfSVSignalStrength],
	))
	s.CarrierPhase = scaleOfSVCarrierPhase * float64(int32(binary.LittleEndian.Uint32(
		b[offset+indexOfSVCarrierPhase:offset+indexOfSVCarrierPhase+lengthOfSVCarrierPhase],
	)))
	s.PseudoRangeResidual = unit.Distance(int32(binary.LittleEndian.Uint32(
		b[offset+indexOfSVPseudoRangeResidual:offset+indexOfSVPseudoRangeResidual+lengthOfSVPseudoRangeResidual],
	))) * unit.Metre
	s.DopplerFrequencyHz = scaleOfSVDopplerFrequency * float64(int32(binary.LittleEndian.Uint32(
		b[offset+indexOfSVDopplerFrequency:offset+indexOfSVDopplerFrequency+lengthOfSVDopplerFrequency],
	)))
	s.Azimuth = scaleOfSVAzimuth * unit.Angle(binary.LittleEndian.Uint16(
		b[offset+indexOfSVAzimuth:offset+indexOfSVAzimuth+lengthOfSVAzimuth],
	)) * unit.Degree
	s.Elevation = scaleOfSVElevation * unit.Angle(binary.LittleEndian.Uint16(
		b[offset+indexOfSVElevation:offset+indexOfSVElevation+lengthOfSVElevation],
	)) * unit.Degree
	return nil
}

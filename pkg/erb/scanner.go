package erb

import (
	"bufio"
	"encoding/binary"
	"io"

	"golang.org/x/xerrors"
)

type Scanner struct {
	sc      *bufio.Scanner
	err     error
	payload []byte
	svIndex int
	ok      bool
	id      ID
	ver     VER
	pos     POS
	stat    STAT
	dops    DOPS
	vel     VEL
	svi     SVI
	sv      SV
}

func NewScanner(r io.Reader) *Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(scanPackets)
	return &Scanner{sc: sc}
}

func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	s.ok = s.sc.Scan()
	s.err = s.sc.Err()
	if !s.ok || s.err != nil {
		return false
	}
	s.id = ID(s.Bytes()[indexOfMessageID])
	lengthOfPayload := binary.LittleEndian.Uint16(
		s.Bytes()[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength],
	)
	s.payload = s.Bytes()[indexOfPayload : indexOfPayload+lengthOfPayload]
	var err error
	switch s.id {
	case IDVER:
		err = s.ver.unmarshal(s.payload)
	case IDPOS:
		err = s.pos.unmarshal(s.payload)
	case IDSTAT:
		err = s.stat.unmarshal(s.payload)
	case IDDOPS:
		err = s.dops.unmarshal(s.payload)
	case IDVEL:
		err = s.vel.unmarshal(s.payload)
	case IDSVI:
		err = s.svi.unmarshal(s.payload)
		s.sv = SV{}
		s.svIndex = 0
	}
	if err != nil {
		s.err = err
		return false
	}
	if s.id == IDVER {
		isProtocolVersionSupported :=
			s.ver.High == SupportedProtocolVersionHigh &&
				s.ver.Medium == SupportedProtocolVersionMedium &&
				s.ver.Low == SupportedProtocolVersionLow
		if !isProtocolVersionSupported {
			s.err = xerrors.Errorf("unsupported protocol version: %d.%d.%d", s.ver.High, s.ver.Medium, s.ver.Low)
			return false
		}
	}
	return true
}

func (s *Scanner) ScanSVI() bool {
	if s.err != nil || s.id != IDSVI || s.svIndex >= int(s.svi.NumSVs) {
		return false
	}
	if err := s.sv.unmarshal(s.payload, s.svIndex); err != nil {
		s.err = err
		return false
	}
	s.svIndex++
	return true
}

func (s *Scanner) ID() ID {
	return s.id
}

func (s *Scanner) VER() *VER {
	return &s.ver
}

func (s *Scanner) POS() *POS {
	return &s.pos
}

func (s *Scanner) STAT() *STAT {
	return &s.stat
}

func (s *Scanner) DOPS() *DOPS {
	return &s.dops
}

func (s *Scanner) VEL() *VEL {
	return &s.vel
}

func (s *Scanner) SVI() *SVI {
	return &s.svi
}

func (s *Scanner) SV() *SV {
	return &s.sv
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) Bytes() []byte {
	return s.sc.Bytes()
}

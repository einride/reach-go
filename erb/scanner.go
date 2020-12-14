package erb

import (
	"bufio"
	"encoding/binary"
	"io"
)

// Scanner provides a convenient interface for reading and parsing ERB messages from a stream.
type Scanner struct {
	sc      *bufio.Scanner
	err     error
	payload []byte
	svIndex int
	id      ID
	ver     VER
	pos     POS
	stat    STAT
	dops    DOPS
	vel     VEL
	svi     SVI
	sv      SV
}

// NewScanner returns a new Scanner to read from r.
func NewScanner(r io.Reader) *Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(ScanPackets)
	return &Scanner{sc: sc}
}

// Scan advances the Scanner to the next message, whose ID will then be
// available through the ID method.
func (c *Scanner) Scan() bool {
	if c.err != nil {
		return false
	}
	if ok := c.sc.Scan(); !ok {
		c.err = c.sc.Err()
		if c.err == nil {
			c.err = io.EOF
		}
		return false
	}
	c.id = ID(c.sc.Bytes()[indexOfMessageID])
	lengthOfPayload := binary.LittleEndian.Uint16(
		c.sc.Bytes()[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength],
	)
	c.payload = c.sc.Bytes()[indexOfPayload : indexOfPayload+lengthOfPayload]
	// assume the scan function has already validated packet types
	switch c.id {
	case IDVER:
		c.ver.unmarshalPayload(c.payload)
	case IDPOS:
		c.pos.unmarshalPayload(c.payload)
	case IDSTAT:
		c.stat.unmarshalPayload(c.payload)
	case IDDOPS:
		c.dops.unmarshalPayload(c.payload)
	case IDVEL:
		c.vel.unmarshalPayload(c.payload)
	case IDSVI:
		c.svi.unmarshalPayload(c.payload)
		c.sv = SV{}
		c.svIndex = 0
	default:
		// allow unknown packets
	}
	return true
}

func (c *Scanner) Err() error {
	if c.err == io.EOF {
		return nil
	}
	return c.err
}

// ScanSVI advances to the next SV packet in an SVI packet.
func (c *Scanner) ScanSVI() bool {
	if c.id != IDSVI || c.svIndex >= int(c.svi.NumSVs) {
		return false
	}
	c.sv.unmarshalPayload(c.payload, c.svIndex)
	c.svIndex++
	return true
}

func (c *Scanner) ID() ID {
	return c.id
}

func (c *Scanner) VER() VER {
	return c.ver
}

func (c *Scanner) POS() POS {
	return c.pos
}

func (c *Scanner) STAT() STAT {
	return c.stat
}

func (c *Scanner) DOPS() DOPS {
	return c.dops
}

func (c *Scanner) VEL() VEL {
	return c.vel
}

func (c *Scanner) SVI() SVI {
	return c.svi
}

func (c *Scanner) SV() SV {
	return c.sv
}

func (c *Scanner) Bytes() []byte {
	return c.sc.Bytes()
}

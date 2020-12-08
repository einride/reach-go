package erb

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// Scanner provides a convenient interface for reading and parsing ERB messages from a stream.
type Scanner struct {
	sc      *bufio.Scanner
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
func (c *Scanner) Scan() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("scan ERB message: %w", err)
		}
	}()
	if ok := c.sc.Scan(); !ok {
		if c.sc.Err() == nil {
			return io.EOF
		}
		return c.sc.Err()
	}
	c.id = ID(c.sc.Bytes()[indexOfMessageID])
	lengthOfPayload := binary.LittleEndian.Uint16(
		c.sc.Bytes()[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength],
	)
	c.payload = c.sc.Bytes()[indexOfPayload : indexOfPayload+lengthOfPayload]
	switch c.id {
	case IDVER:
		err = c.ver.unmarshal(c.payload)
	case IDPOS:
		err = c.pos.unmarshal(c.payload)
	case IDSTAT:
		err = c.stat.unmarshal(c.payload)
	case IDDOPS:
		err = c.dops.unmarshal(c.payload)
	case IDVEL:
		err = c.vel.unmarshal(c.payload)
	case IDSVI:
		err = c.svi.unmarshal(c.payload)
		c.sv = SV{}
		c.svIndex = 0
	}
	if err != nil {
		return err
	}
	if c.id == IDVER {
		isProtocolVersionSupported :=
			c.ver.High == SupportedProtocolVersionHigh &&
				c.ver.Medium == SupportedProtocolVersionMedium &&
				c.ver.Low == SupportedProtocolVersionLow
		if !isProtocolVersionSupported {
			return fmt.Errorf(
				"unsupported protocol version: %d.%d.%d",
				c.ver.High,
				c.ver.Medium,
				c.ver.Low,
			)
		}
	}
	return nil
}

// ScanSVI advances to the next SV packet in an SVI packet.
func (c *Scanner) ScanSVI() bool {
	if c.id != IDSVI || c.svIndex >= int(c.svi.NumSVs) {
		return false
	}
	if err := c.sv.unmarshal(c.payload, c.svIndex); err != nil {
		// this should not happen if our scan function is implemented properly
		return false
	}
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

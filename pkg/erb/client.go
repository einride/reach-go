package erb

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"time"

	"golang.org/x/xerrors"
)

// Client is an ERB client.
type Client struct {
	sc      *bufio.Scanner
	conn    Conn
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

// Conn is an interface for the connection used by the ERB client.
type Conn interface {
	io.ReadCloser
	SetReadDeadline(time.Time) error
}

// NewClient creates a new ERB client with the provided connection.
func NewClient(conn Conn) *Client {
	sc := bufio.NewScanner(conn)
	sc.Split(scanPackets)
	return &Client{sc: sc, conn: conn}
}

// Receive an ERB message on the connection.
func (c *Client) Receive(ctx context.Context) error {
	deadline, _ := ctx.Deadline()
	if err := c.conn.SetReadDeadline(deadline); err != nil {
		return xerrors.Errorf("erb client: receive: %w", err)
	}
	if ok := c.sc.Scan(); !ok {
		if c.sc.Err() == nil {
			return xerrors.Errorf("erb client: receive: %w", io.EOF)
		}
		return xerrors.Errorf("erb client: receive: %w", c.sc.Err())
	}
	c.id = ID(c.sc.Bytes()[indexOfMessageID])
	lengthOfPayload := binary.LittleEndian.Uint16(
		c.sc.Bytes()[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength],
	)
	c.payload = c.sc.Bytes()[indexOfPayload : indexOfPayload+lengthOfPayload]
	var err error
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
		return xerrors.Errorf("erb client: receive: %w", err)
	}
	if c.id == IDVER {
		isProtocolVersionSupported :=
			c.ver.High == SupportedProtocolVersionHigh &&
				c.ver.Medium == SupportedProtocolVersionMedium &&
				c.ver.Low == SupportedProtocolVersionLow
		if !isProtocolVersionSupported {
			return xerrors.Errorf(
				"erb client: receive: unsupported protocol version: %d.%d.%d",
				c.ver.High, c.ver.Medium, c.ver.Low,
			)
		}
	}
	return nil
}

func (c *Client) ScanSVI() bool {
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

func (c *Client) ID() ID {
	return c.id
}

func (c *Client) VER() *VER {
	return &c.ver
}

func (c *Client) POS() *POS {
	return &c.pos
}

func (c *Client) STAT() *STAT {
	return &c.stat
}

func (c *Client) DOPS() *DOPS {
	return &c.dops
}

func (c *Client) VEL() *VEL {
	return &c.vel
}

func (c *Client) SVI() *SVI {
	return &c.svi
}

func (c *Client) SV() *SV {
	return &c.sv
}

func (c *Client) Bytes() []byte {
	return c.sc.Bytes()
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return xerrors.Errorf("erb client: close: %w", err)
	}
	return nil
}

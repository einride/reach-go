package erb

// ID represents an ERB message ID.
type ID uint8

//go:generate stringer -type=ID -trimprefix=ID

const (
	// IDVER is the ID of the VER message.
	IDVER ID = 0x01
	// IDPOS is the ID of the POS message.
	IDPOS ID = 0x02
	// IDSTAT is the ID of the STAT message.
	IDSTAT ID = 0x03
	// IDDOPS is the ID of the DOPS message.
	IDDOPS ID = 0x04
	// IDVEL is the ID of the VEL message.
	IDVEL ID = 0x05
	// IDSVI is the ID of the SVI message.
	IDSVI ID = 0x06
)

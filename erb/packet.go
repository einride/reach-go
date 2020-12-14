package erb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	syncChar1 = 'E'
	syncChar2 = 'R'
)

const (
	indexOfSyncWord       = 0
	lengthOfSyncWord      = 2
	indexOfMessageID      = indexOfSyncWord + lengthOfSyncWord
	lengthOfMessageID     = 1
	indexOfPayloadLength  = indexOfMessageID + lengthOfMessageID
	lengthOfPayloadLength = 2
	indexOfPayload        = indexOfPayloadLength + lengthOfPayloadLength
	lengthOfChecksum      = 2
)

func calculateIndexOfChecksum(lengthOfPayload uint16) int {
	return indexOfPayload + int(lengthOfPayload)
}

func calculateLengthOfPacket(lengthOfPayload uint16) int {
	return indexOfPayload + int(lengthOfPayload) + lengthOfChecksum
}

// ScanPackets is a split function for a bufio.Scanner that returns each ERB packet.
func ScanPackets(data []byte, _ bool) (advance int, token []byte, err error) {
	if len(data) < indexOfPayloadLength+lengthOfPayloadLength {
		return 0, nil, nil
	}
	// scan until start of packet
	if data[0] != syncChar1 || data[1] != syncChar2 {
		i := bytes.Index(data, []byte{syncChar1, syncChar2})
		if i == -1 {
			return len(data), nil, nil
		}
		return i, nil, nil
	}
	// parse payload length
	lengthOfPayload := binary.LittleEndian.Uint16(
		data[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength],
	)
	lengthOfPacket := calculateLengthOfPacket(lengthOfPayload)
	if len(data) < lengthOfPacket {
		return 0, nil, nil
	}
	packet := data[:lengthOfPacket]
	payload := data[indexOfPayload : indexOfPayload+lengthOfPayload]
	messageID := ID(packet[indexOfMessageID])
	// verify checksum
	expectedChecksum := fletcher(packet[indexOfMessageID:calculateIndexOfChecksum(lengthOfPayload)])
	actualChecksum := binary.LittleEndian.Uint16(packet[calculateIndexOfChecksum(lengthOfPayload):])
	if expectedChecksum != actualChecksum {
		return 0, nil, fmt.Errorf(
			"checksum mismatch (expected 0x%x but got 0x%x)",
			expectedChecksum,
			actualChecksum,
		)
	}
	if err := validatePayload(messageID, payload); err != nil {
		return 0, nil, err
	}
	return lengthOfPacket, packet, nil
}

func fletcher(data []byte) uint16 {
	var a, b uint8
	for i := 0; i < len(data); i++ {
		a += data[i]
		b += a
	}
	return uint16(b)<<8 | uint16(a)
}

func validatePayload(messageID ID, payload []byte) error {
	switch messageID {
	case IDVER:
		return validatePayloadVER(payload)
	case IDPOS:
		return validatePayloadPOS(payload)
	case IDSTAT:
		return validatePayloadSTAT(payload)
	case IDDOPS:
		return validatePayloadDOPS(payload)
	case IDVEL:
		return validatePayloadVEL(payload)
	case IDSVI:
		return validatePayloadSVI(payload)
	default:
		return nil // allow unknown packets
	}
}

func validatePayloadVER(payload []byte) error {
	if len(payload) != lengthOfVER {
		return fmt.Errorf("validate VER payload: illegal length %d (expected %d)", len(payload), lengthOfVER)
	}
	var ver VER
	ver.unmarshalPayload(payload)
	isProtocolVersionSupported :=
		ver.High == SupportedProtocolVersionHigh &&
			ver.Medium == SupportedProtocolVersionMedium &&
			ver.Low == SupportedProtocolVersionLow
	if !isProtocolVersionSupported {
		return fmt.Errorf(
			"validate VER payload: unsupported protocol version %d.%d.%d",
			ver.High,
			ver.Medium,
			ver.Low,
		)
	}
	return nil
}

func validatePayloadPOS(payload []byte) error {
	if len(payload) != lengthOfPOS {
		return fmt.Errorf("validate %v payload: illegal length %d (expected %d)", IDPOS, len(payload), lengthOfPOS)
	}
	return nil
}

func validatePayloadSTAT(payload []byte) error {
	if len(payload) != lengthOfSTAT {
		return fmt.Errorf("validate %v payload: illegal length %d (expected %d)", IDSTAT, len(payload), lengthOfSTAT)
	}
	return nil
}

func validatePayloadDOPS(payload []byte) error {
	if len(payload) != lengthOfDOPS {
		return fmt.Errorf("validate %v payload: illegal length %d (expected %d)", IDDOPS, len(payload), lengthOfDOPS)
	}
	return nil
}

func validatePayloadVEL(payload []byte) error {
	if len(payload) != lengthOfVEL {
		return fmt.Errorf("validate %v payload: illegal length %d (expected %d)", IDVEL, len(payload), lengthOfVEL)
	}
	return nil
}

func validatePayloadSVI(payload []byte) error {
	const minLength = indexOfNumSVs + lengthOfNumSVs
	if len(payload) < minLength {
		return fmt.Errorf(
			"validate %v payload: illegal length %d (expected minimum %d)", IDSVI, len(payload), lengthOfSVI,
		)
	}
	var svi SVI
	svi.unmarshalPayload(payload)
	expectedLength := minLength + lengthOfSV*int(svi.NumSVs)
	if len(payload) != expectedLength {
		return fmt.Errorf("validate %v payload: illegal length %d (expected %d)", IDSVI, len(payload), expectedLength)
	}
	return nil
}

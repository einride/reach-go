package erb

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/xerrors"
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

func indexOfChecksum(lengthOfPayload uint16) int {
	return indexOfPayload + int(lengthOfPayload)
}

func lengthOfPacket(lengthOfPayload uint16) int {
	return indexOfPayload + int(lengthOfPayload) + lengthOfChecksum
}

func scanPackets(data []byte, _ bool) (advance int, token []byte, err error) {
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
	payloadLength := binary.LittleEndian.Uint16(data[indexOfPayloadLength : indexOfPayloadLength+lengthOfPayloadLength])
	packetLength := lengthOfPacket(payloadLength)
	if len(data) < packetLength {
		return 0, nil, nil
	}
	packet := data[:packetLength]
	// verify checksum
	expectedChecksum := fletcher(packet[indexOfMessageID:indexOfChecksum(payloadLength)])
	actualChecksum := binary.LittleEndian.Uint16(packet[indexOfChecksum(payloadLength):])
	if expectedChecksum != actualChecksum {
		return 0, nil, xerrors.Errorf(
			"checksum mismatch: expected: 0x%x, actual: 0x%x",
			expectedChecksum,
			actualChecksum,
		)
	}
	return packetLength, packet, nil
}

func fletcher(data []byte) uint16 {
	var a, b uint8
	for i := 0; i < len(data); i++ {
		a += data[i]
		b += a
	}
	return uint16(b)<<8 | uint16(a)
}

package udp

import bytehelpers "networking/internal/byte_helpers"

var udpProtocolNumber = 17

type PseudoHeader struct {
	SourceIP [4]byte
	DestIP   [4]byte
	Protocol uint8
	Length   uint16
}

func NewPseudoHeader(source, dest [4]byte, udpLength uint16) PseudoHeader {
	return PseudoHeader{
		Protocol: uint8(udpProtocolNumber),
		SourceIP: source,
		DestIP:   dest,
		Length:   udpLength,
	}
}

func (p *PseudoHeader) CreatePseudoHeader() []byte {
	if p == nil {
		return []byte{}
	}

	sourceIPBytes := p.SourceIP[:]
	destIPBytes := p.DestIP[:]
	lengthBytes := bytehelpers.Uint16ToByteArray(p.Length)
	zeroByte := []byte{0}

	return bytehelpers.ConcatenateByteArrays(
		sourceIPBytes,
		destIPBytes,
		zeroByte,
		[]byte{p.Protocol},
		lengthBytes,
	)
}

func (p PseudoHeader) ParsePseudoHeader(data []byte) PseudoHeader {
	sourceIP := [4]byte{data[0], data[1], data[2], data[3]}
	destIP := [4]byte{data[4], data[5], data[6], data[7]}
	length := bytehelpers.ByteArrayToUint16([]byte{data[10], data[11]})

	return PseudoHeader{
		SourceIP: sourceIP,
		DestIP:   destIP,
		Length:   length,
		Protocol: uint8(udpProtocolNumber),
	}
}

package udp

import (
	"context"
	"fmt"

	"networking/internal/byte_helpers"
	"networking/internal/logger"
)

// UDPGram represents a UDP datagram structure
type UDPGram struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
	Data            []byte
}

// NewUDPGram Helper function to create a new UDP datagram struct.
// Used for internal code and consumers of the package
func (h *UDPGram) NewUDPGram(ctx *context.Context, sourcePort, destinationPort *uint16, data *[]byte) (*UDPGram, error) {
	logger := logger.GetLoggerFromContext(*ctx, nil)

	if sourcePort == nil || *sourcePort == 0 {
		logger.Warn("Source port is set to 0")
	}

	if destinationPort == nil || *destinationPort == 0 {
		err := fmt.Errorf("invalid destination port value: %d", *destinationPort)
		logger.Error(err.Error())

		return nil, err
	}

	if data == nil {
		err := fmt.Errorf("data cannot be nil")
		logger.Error(err.Error())

		return nil, err
	}

	if len(*data) == 0 {
		err := fmt.Errorf("data cannot be empty")
		logger.Error(err.Error())

		return nil, err
	}

	length := uint16(8 + len(*data))
	checksum := bytehelpers.CreateOnesComplementChecksum(*data)

	return &UDPGram{
		SourcePort:      *sourcePort,
		DestinationPort: *destinationPort,
		Length:          length,
		Checksum:        checksum,
		Data:            *data,
	}, nil
}

// CreateUDPGram Function to create a raw UDP datagram byte array from the UDPGram struct
func (h *UDPGram) CreateUDPGram(ctx *context.Context) ([]byte, error) {
	logger := logger.GetLoggerFromContext(*ctx, nil)

	if h.DestinationPort == 0 {
		err := fmt.Errorf("invalid destination port value: %d", h.DestinationPort)
		logger.Error(err.Error())

		return nil, err
	}

	if len(h.Data) == 0 {
		err := fmt.Errorf("data cannot be empty")
		logger.Error(err.Error())

		return nil, err
	}

	if h.SourcePort == 0 {
		logger.Warn("Source port is set to 0")
	}

	sourcePortBytes := bytehelpers.Uint16ToByteArray(h.SourcePort)

	destinationPortBytes := bytehelpers.Uint16ToByteArray(h.DestinationPort)

	length := uint16(8 + len(h.Data))
	lengthBytes := bytehelpers.Uint16ToByteArray(length)

	checksum := bytehelpers.CreateOnesComplementChecksum(h.Data)
	checksumBytes := bytehelpers.Uint16ToByteArray(checksum)

	header := bytehelpers.ConcatenateByteArrays(sourcePortBytes, destinationPortBytes, lengthBytes, checksumBytes, h.Data)

	return header, nil
}

// ParseRawUDPGram Function to parse a raw UDP datagram byte array into a UDPGram struct
func ParseRawUDPGram(ctx context.Context, data []byte) (*UDPGram, error) {
	logger := logger.GetLoggerFromContext(ctx, nil)

	sourcePort := bytehelpers.ByteArrayToUint16(data[0:2])
	destinationPort := bytehelpers.ByteArrayToUint16(data[2:4])
	length := bytehelpers.ByteArrayToUint16(data[4:6])
	checksum := bytehelpers.ByteArrayToUint16(data[6:8])

	if len(data) != int(length) {
		err := fmt.Errorf("invalid UDP header length. Expected length (%d) does not match actual data length (%d)", length, len(data))
		logger.Error(err.Error())

		return nil, err
	}

	if checksum == 0 {
		logger.Warn("Checksum is set to 0")
	}

	if sourcePort == 0 {
		logger.Warn("Source port is set to 0")
	}

	if destinationPort == 0 {
		err := fmt.Errorf("invalid destination port value: %d", destinationPort)
		logger.Error(err.Error())

		return nil, err
	}

	return &UDPGram{
		SourcePort:      sourcePort,
		DestinationPort: destinationPort,
		Length:          length,
		Checksum:        checksum,
		Data:            data[8:],
	}, nil
}

func (h *UDPGram) IsEqual(a *UDPGram) bool {
	if a == nil {
		return false
	}

	areSourcePortsEqual := h.SourcePort == a.SourcePort
	areDestinationPortsEqual := h.DestinationPort == a.DestinationPort
	areLengthsEqual := h.Length == a.Length
	areChecksumsEqual := h.Checksum == a.Checksum
	areDataEqual := bytehelpers.AreByteArraysEqual(h.Data, a.Data)

	return areSourcePortsEqual && areDestinationPortsEqual && areLengthsEqual && areChecksumsEqual && areDataEqual
}

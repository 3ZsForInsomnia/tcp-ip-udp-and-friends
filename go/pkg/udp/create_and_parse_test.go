package udp

import (
	"context"
	"networking/internal/logger"
	testhelpers "networking/internal/test_helpers"
	"testing"
)

/**
* Test cases for Parsing UDP datagrams
 */
func Test_Parse_HappyPath(t *testing.T) {
	lctx := logger.PrepTest()

	data := []byte("Hello, this is a test UDP datagram.")
	expected := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Length:          uint16(8 + len(data)),
		Checksum:        0x1A2B,
		Data:            data,
	}

	// Source Port: 8080 (0x1F90), Dest: 80 (0x0050), Length: 43 (0x002B), Checksum: 0x1A2B
	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x2B, // Length: 43
		0x1A, 0x2B, // Checksum
	}
	input = append(input, data...)

	actual, err := ParseRawUDPGram(*lctx, input)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if !expected.IsEqual(actual) {
		t.Errorf("Parsed UDPGram does not match expected.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

func Test_Parse_InvalidLength(t *testing.T) {
	lctx := logger.PrepTest()

	expected := "invalid UDP header length. Expected length (100) does not match actual data length (13)"
	data := []byte("Short")
	// Length field says 100 but actual data is only 13 bytes total
	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x64, // Length: 100 (incorrect)
		0x1A, 0x2B, // Checksum
	}
	input = append(input, data...)

	_, err := ParseRawUDPGram(*lctx, input)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s', but got '%v'", expected, err)
	}
}

func Test_Parse_ChecksumZero(t *testing.T) {
	data := []byte("Checksum zero")
	expected := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Length:          uint16(8 + len(data)),
		Checksum:        0,
		Data:            data,
	}

	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x15, // Length: 21
		0x00, 0x00, // Checksum: 0
	}
	input = append(input, data...)

	actual, err := ParseRawUDPGram(context.TODO(), input)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if !expected.IsEqual(actual) {
		t.Errorf("Parsed UDPGram does not match expected.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

func Test_Parse_IncorrectChecksum(t *testing.T) {
	data := []byte("Bad checksum")
	expected := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Length:          uint16(8 + len(data)),
		Checksum:        0xFFFF,
		Data:            data,
	}

	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x14, // Length: 20
		0xFF, 0xFF, // Checksum: 0xFFFF (wrong)
	}
	input = append(input, data...)

	actual, err := ParseRawUDPGram(context.TODO(), input)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if !expected.IsEqual(actual) {
		t.Errorf("Parsed UDPGram does not match expected.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

func Test_Parse_InvalidDestinationPort(t *testing.T) {
	expected := "invalid destination port value: 0"

	data := []byte("Test")
	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x00, // Destination Port: 0 (invalid)
		0x00, 0x0C, // Length: 12
		0x1A, 0x2B, // Checksum
	}
	input = append(input, data...)

	_, err := ParseRawUDPGram(context.TODO(), input)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s', but got '%v'", expected, err)
	}
}

func Test_Parse_SourcePortZero(t *testing.T) {
	data := []byte("Source zero")
	expected := UDPGram{
		SourcePort:      0,
		DestinationPort: 80,
		Length:          uint16(8 + len(data)),
		Checksum:        0x1A2B,
		Data:            data,
	}

	input := []byte{
		0x00, 0x00, // Source Port: 0
		0x00, 0x50, // Destination Port: 80
		0x00, 0x13, // Length: 19
		0x1A, 0x2B, // Checksum
	}
	input = append(input, data...)

	actual, err := ParseRawUDPGram(context.TODO(), input)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if !expected.IsEqual(actual) {
		t.Errorf("Parsed UDPGram does not match expected.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

func Test_Parse_NoData(t *testing.T) {
	expected := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Length:          8,
		Checksum:        0x1A2B,
		Data:            []byte{},
	}

	input := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x08, // Length: 8 (header only)
		0x1A, 0x2B, // Checksum
	}

	actual, err := ParseRawUDPGram(context.TODO(), input)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if !expected.IsEqual(actual) {
		t.Errorf("Parsed UDPGram does not match expected.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

/**
* Test cases for Creating UDP datagrams
 */

func Test_Create_HappyPath(t *testing.T) {
	lctx := logger.PrepTest()
	udpGram := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Data:            []byte("Hello UDP"),
	}

	expected := []byte{
		0x1F, 0x90, // Source Port: 8080
		0x00, 0x50, // Destination Port: 80
		0x00, 0x11, // Length: 17 (8 + 9)
	}

	actual, err := udpGram.CreateUDPGram(lctx)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if len(actual) < len(expected) || string(actual[0:6]) != string(expected) {
		t.Errorf("Created UDPGram does not match expected header.\nExpected: %+v\nActual: %+v", expected, actual[0:6])
	}
}

func Test_Create_InvalidDestinationPort(t *testing.T) {
	lctx := logger.PrepTest()

	udpGram := UDPGram{
		SourcePort:      8080,
		DestinationPort: 0,
		Data:            []byte("Test"),
	}

	expected := "invalid destination port value: 0"

	_, err := udpGram.CreateUDPGram(lctx)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s', but got '%v'", expected, err)
	}
}

func Test_Create_SourcePortZero(t *testing.T) {
	lctx := logger.PrepTest()

	udpGram := UDPGram{
		SourcePort:      0,
		DestinationPort: 80,
		Data:            []byte("Test"),
	}

	expected := []byte{
		0x00, 0x00, // Source Port: 0
		0x00, 0x50, // Destination Port: 80
		0x00, 0x0C, // Length: 12 (8 + 4)
	}

	actual, err := udpGram.CreateUDPGram(lctx)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	if len(actual) < len(expected) || string(actual[0:6]) != string(expected) {
		t.Errorf("Created UDPGram does not match expected header.\nExpected: %+v\nActual: %+v", expected, actual[0:6])
	}
}

func Test_Create_NoData(t *testing.T) {
	lctx := logger.PrepTest()

	udpGram := UDPGram{
		SourcePort:      8080,
		DestinationPort: 80,
		Data:            []byte{},
	}

	expected := "data cannot be empty"

	_, err := udpGram.CreateUDPGram(lctx)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s', but got '%v'", expected, err)
	}
}

func Test_Create_NilSourcePort(t *testing.T) {
	lctx := logger.PrepTest()

	udpGram := UDPGram{
		SourcePort:      0,
		DestinationPort: 80,
		Data:            []byte("Test"),
	}

	expectedWarning := "Source port is set to 0"

	_, err := udpGram.CreateUDPGram(lctx)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	mockLogger := logger.GetLoggerFromContext(*lctx, nil).(*logger.MockLogger)
	if !mockLogger.HasWarning(expectedWarning) {
		t.Errorf("Expected warning '%s' not found in logs", expectedWarning)
	}
}

func Test_Create_NilDestinationPort(t *testing.T) {
	lctx := logger.PrepTest()

	udpGram := UDPGram{
		SourcePort:      8080,
		DestinationPort: 0,
		Data:            []byte("Test"),
	}

	expected := "invalid destination port value: 0"

	_, err := udpGram.CreateUDPGram(lctx)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected error '%s', but got '%v'", expected, err)
	}
}

/**
* End-to-End Test cases for Creating and Parsing UDP datagrams
 */

func Test_Create_Then_Parse_Consistency(t *testing.T) {
	lctx := logger.PrepTest()

	original := UDPGram{
		SourcePort:      12345,
		DestinationPort: 443,
		Data:            []byte("Round-trip test"),
	}

	// Create the raw UDP datagram
	rawBytes, err := original.CreateUDPGram(lctx)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	// Parse it back
	parsed, err := ParseRawUDPGram(*lctx, rawBytes)
	testhelpers.FailTestIfErrorIsPresent(t, err)

	// Verify the parsed datagram matches the original
	if parsed.SourcePort != original.SourcePort {
		t.Errorf("Source port mismatch. Expected: %d, Got: %d", original.SourcePort, parsed.SourcePort)
	}

	if parsed.DestinationPort != original.DestinationPort {
		t.Errorf("Destination port mismatch. Expected: %d, Got: %d", original.DestinationPort, parsed.DestinationPort)
	}

	if string(parsed.Data) != string(original.Data) {
		t.Errorf("Data mismatch. Expected: %s, Got: %s", original.Data, parsed.Data)
	}
}

package tftp_test

import (
	"fmt"
	"networking/pkg/tftp"
	"testing"
)

/**
* Test cases for TFTP Ack packet serialization and deserialization
 */

func TestSerializeAckPacket(t *testing.T) {
	expected := []byte{0, 4, 0, 1} // Opcode 4 (ACK) and Block Number 1
	blockNum := uint16(1)
	inputPacket := &tftp.TFTPPacket{
		Opcode:   4,
		BlockNum: &blockNum,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}

	if string(actual) != string(expected) {
		t.Fatalf("Expected bytes %v, got %v", expected, actual)
	}
}

func TestSerializeAckPacketNoBlockNum(t *testing.T) {
	expectedErrMsg := "cannot serialize ACK packet: Block number is nil"
	inputPacket := &tftp.TFTPPacket{
		Opcode:   4,
		BlockNum: nil,
	}

	_, err := inputPacket.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to missing BlockNum, but got none")
	}

	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestDeserializeAckPacket(t *testing.T) {
	expectedBlockNum := uint16(1)
	expectedData := []byte("Hi there this is some heckin' cool data")

	expectedPacket := &tftp.TFTPPacket{
		BlockNum: &expectedBlockNum,
		Data:     &expectedData,
		Opcode:   3,
	}

	input := []byte{0, 3, 0, 1, 'H', 'i', ' ', 't', 'h', 'e', 'r', 'e', ' ', 't', 'h', 'i', 's', ' ', 'i', 's', ' ', 's', 'o', 'm', 'e', ' ', 'h', 'e', 'c', 'k', 'i', 'n', '\'', ' ', 'c', 'o', 'o', 'l', ' ', 'd', 'a', 't', 'a'}

	actualPacket, err := tftp.Deserialize(input)

	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}

	if actualPacket.Opcode != expectedPacket.Opcode {
		t.Fatalf("Expected Opcode %d, got %d", expectedPacket.Opcode, actualPacket.Opcode)
	} else if *actualPacket.BlockNum != *expectedPacket.BlockNum {
		t.Fatalf("Expected BlockNum %d, got %d", *expectedPacket.BlockNum, *actualPacket.BlockNum)
	} else if string(*actualPacket.Data) != string(*expectedPacket.Data) {
		t.Fatalf("Expected Data %s, got %s", string(*expectedPacket.Data), string(*actualPacket.Data))
	}
}

func TestDeserializeAckPacketNoData(t *testing.T) {
	expectedErrMsg := "ACK packet too short to contain Block Number"
	input := []byte{0, 4} // Incomplete ACK packet

	_, err := tftp.Deserialize(input)

	if err == nil {
		t.Fatalf("Expected error due to incomplete ACK packet, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

/**
* Test cases for TFTP Error packet serialization and deserialization
 */

func TestSerializeErrorPacket(t *testing.T) {
	expected := []byte{0, 5, 0, 1, 'F', 'i', 'l', 'e', ' ', 'n', 'o', 't', ' ', 'f', 'o', 'u', 'n', 'd', '.', 0} // Opcode 5 (ERROR), Error Code 1, "File not found"

	errCode := uint16(1)
	inputPacket := &tftp.TFTPPacket{
		ErrorCode: &errCode,
		Opcode:    5,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}

	if string(actual) != string(expected) {
		t.Fatalf("Expected bytes %v, got %v", expected, actual)
	}
}

func TestSerializeErrorPacketNoErrorCode(t *testing.T) {
	expectedErrMsg := "cannot serialize ERROR packet: ErrorCode is nil"
	inputPacket := &tftp.TFTPPacket{
		Opcode:    5,
		ErrorCode: nil,
	}

	_, err := inputPacket.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to missing ErrorCode, but got none")
	}

	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSerializeErrorPacketUnknownErrorCode(t *testing.T) {
	expectedErrMsg := "unknown error code: 9999"
	inputErrCode := uint16(9999)
	packet := &tftp.TFTPPacket{
		Opcode:    5,
		ErrorCode: &inputErrCode,
	}

	_, err := packet.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to unknown ErrorCode, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestDeserializeErrorPacket(t *testing.T) {
	expectedErrCode := uint16(1)
	expectedErrMsg := "File not found."

	input := []byte{0, 5, 0, 1, 'F', 'i', 'l', 'e', ' ', 'n', 'o', 't', ' ', 'f', 'o', 'u', 'n', 'd', '.', 0}

	actualPacket, err := tftp.Deserialize(input)
	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}

	if actualPacket.Opcode != 5 {
		t.Fatalf("Expected Opcode 5, got %d", actualPacket.Opcode)
	}
	if *actualPacket.ErrorCode != expectedErrCode {
		t.Fatalf("Expected ErrorCode %d, got %d", expectedErrCode, *actualPacket.ErrorCode)
	}
	if *actualPacket.ErrorMsg != expectedErrMsg {
		t.Fatalf("Expected ErrorMsg '%s', got '%s'", expectedErrMsg, *actualPacket.ErrorMsg)
	}
}

func TestDeserializeErrorPacketNoErrorCode(t *testing.T) {
	expectedErrMsg := "ERROR packet too short to contain Error Code"
	input := []byte{0, 5} // Incomplete ERROR packet

	_, err := tftp.Deserialize(input)
	if err == nil {
		t.Fatalf("Expected error due to incomplete ERROR packet, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

/**
* Test cases for TFTP Data packet serialization and deserialization
 */

func TestSerializeDataPacket(t *testing.T) {
	expected := []byte{0, 3, 0, 1, 'H', 'i', ' ', 't', 'h', 'e', 'r', 'e', ' ', 't', 'h', 'i', 's', ' ', 'i', 's', ' ', 's', 'o', 'm', 'e', ' ', 'h', 'e', 'c', 'k', 'i', 'n', '\'', ' ', 'c', 'o', 'o', 'l', ' ', 'd', 'a', 't', 'a'} // Opcode 3 (DATA), Block Number 1, Data

	inputData := []byte("Hi there this is some heckin' cool data")
	inputBlockNum := uint16(1)
	inputPacket := &tftp.TFTPPacket{
		Data:     &inputData,
		BlockNum: &inputBlockNum,
		Opcode:   3,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}

	if string(actual) != string(expected) {
		t.Fatalf("Expected bytes %v, got %v", expected, actual)
	}
}

func TestSerializeDataPacketNoBlockNum(t *testing.T) {
	expectedErrMsg := "cannot serialize DATA packet: Block number is nil"

	inputData := []byte("Some data")
	inputPacket := &tftp.TFTPPacket{
		Data:   &inputData,
		Opcode: 3,
	}

	_, err := inputPacket.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to missing BlockNum, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSerializeDataPacketNoData(t *testing.T) {
	inputBlockNum := uint16(1)
	inputPacket := &tftp.TFTPPacket{
		Data:     nil,
		BlockNum: &inputBlockNum,
		Opcode:   3,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}
	if string(actual) != string([]byte{0, 3, 0, 1}) {
		t.Fatalf("Expected bytes %v, got %v", []byte{0, 3, 0, 1}, actual)
	}
}

func TestDeserializeDataPacket(t *testing.T) {
	expectedBlockNum := uint16(1)
	expectedData := []byte("Hi there this is some heckin' cool data")

	input := []byte{0, 3, 0, 1, 'H', 'i', ' ', 't', 'h', 'e', 'r', 'e', ' ', 't', 'h', 'i', 's', ' ', 'i', 's', ' ', 's', 'o', 'm', 'e', ' ', 'h', 'e', 'c', 'k', 'i', 'n', '\'', ' ', 'c', 'o', 'o', 'l', ' ', 'd', 'a', 't', 'a'}

	actualPacket, err := tftp.Deserialize(input)
	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}

	if actualPacket.Opcode != 3 {
		t.Fatalf("Expected Opcode 3, got %d", actualPacket.Opcode)
	}
	if *actualPacket.BlockNum != expectedBlockNum {
		t.Fatalf("Expected BlockNum %d, got %d", expectedBlockNum, *actualPacket.BlockNum)
	}
	if string(*actualPacket.Data) != string(expectedData) {
		t.Fatalf("Expected Data '%s', got '%s'", string(expectedData), string(*actualPacket.Data))
	}
}

func TestDeserializeDataPacketNoData(t *testing.T) {
	expectedBlockNum := uint16(1)
	expectedData := []byte{}

	input := []byte{0, 3, 0, 1} // DATA packet with no data

	actualPacket, err := tftp.Deserialize(input)
	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}
	if actualPacket.Opcode != 3 {
		t.Fatalf("Expected Opcode 3, got %d", actualPacket.Opcode)
	}
	if *actualPacket.BlockNum != expectedBlockNum {
		t.Fatalf("Expected BlockNum %d, got %d", expectedBlockNum, *actualPacket.BlockNum)
	}
	if string(*actualPacket.Data) != string(expectedData) {
		t.Fatalf("Expected Data '%s', got '%s'", string(expectedData), string(*actualPacket.Data))
	}
}

func TestDeserializeDataPacketNoBlockNum(t *testing.T) {
	expectedErrMsg := "DATA packet too short to contain Block Number"

	input := []byte{0, 3} // Incomplete DATA packet

	_, err := tftp.Deserialize(input)
	if err == nil {
		t.Fatalf("Expected error due to incomplete DATA packet, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

/**
* Test cases for TFTP RRQ packet serialization and deserialization
 */

func TestSerializeRRQPacket(t *testing.T) {
	expected := []byte{0, 1, 'f', 'i', 'l', 'e', '.', 't', 'x', 't', 0, 'o', 'c', 't', 'e', 't', 0} // Opcode 1 (RRQ), Filename "file.txt", Mode "octet"

	inputMode, err := tftp.StringToMode("octet")
	if err != nil {
		t.Fatalf("Unexpected error converting mode string to Mode: %v", err)
	}

	inputFilename := "file.txt"

	connection := &tftp.TFTPConnection{
		Filename: &inputFilename,
		Mode:     &inputMode,
	}

	inputPacket := &tftp.TFTPPacket{
		Connection: *connection,
		Opcode:     1,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}
	if string(actual) != string(expected) {
		t.Fatalf("Expected bytes %v, got %v", expected, actual)
	}
}

func TestSerializeRRQPacketNoFilename(t *testing.T) {
	expectedErrMsg := "filename is nil or empty"

	inputMode, err := tftp.StringToMode("octet")
	if err != nil {
		t.Fatalf("Unexpected error converting mode string to Mode: %v", err)
	}

	connection := &tftp.TFTPConnection{
		Filename: nil,
		Mode:     &inputMode,
	}

	inputPacket := &tftp.TFTPPacket{
		Connection: *connection,
		Opcode:     1,
	}

	_, err = inputPacket.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to missing Filename, but got none")
	}

	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSerializeRRQPacketNoMode(t *testing.T) {
	expectedErrMsg := "mode is nil"

	inputFilename := "file.txt"

	connection := &tftp.TFTPConnection{
		Filename: &inputFilename,
		Mode:     nil,
	}

	inputPacket := &tftp.TFTPPacket{
		Connection: *connection,
		Opcode:     1,
	}

	_, err := inputPacket.Serialize()
	if err == nil {
		t.Fatalf("Expected error due to missing Mode, but got none")
	}

	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestDeserializeRRQPacket(t *testing.T) {
	expectedFilename := "file.txt"
	expectedMode, err := tftp.StringToMode("octet")
	fmt.Println("Expected Mode:", expectedMode)
	if err != nil {
		t.Fatalf("Unexpected error converting mode string to Mode: %v", err)
	}

	input := []byte{0, 1, 'f', 'i', 'l', 'e', '.', 't', 'x', 't', 0, 'o', 'c', 't', 'e', 't', 0}

	actualPacket, err := tftp.Deserialize(input)
	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}

	if actualPacket.Opcode != 1 {
		t.Fatalf("Expected Opcode 1, got %d", actualPacket.Opcode)
	}
	if *actualPacket.Connection.Filename != expectedFilename {
		t.Fatalf("Expected Filename '%s', got '%s'", expectedFilename, *actualPacket.Connection.Filename)
	}
	if *actualPacket.Connection.Mode != expectedMode {
		t.Fatalf("Expected Mode %d, got %d", expectedMode, *actualPacket.Connection.Mode)
	}
}

func TestDeserializeRRQPacketNoFilename(t *testing.T) {
	expectedErrMsg := "mode string is empty"

	input := []byte{0, 1, 0} // RRQ packet with empty filename

	_, err := tftp.Deserialize(input)
	if err == nil {
		t.Fatalf("Expected error due to missing Filename, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestDeserializeRRQPacketNoMode(t *testing.T) {
	expectedErrMsg := "mode string is empty"

	input := []byte{0, 1, 'f', 'i', 'l', 'e', '.', 't', 'x', 't', 0, 0} // RRQ packet with empty mode

	_, err := tftp.Deserialize(input)
	if err == nil {
		t.Fatalf("Expected error due to missing Mode, but got none")
	}
	if err.Error() != expectedErrMsg {
		t.Fatalf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

/**
* Test cases for TFTP WRQ packet serialization and deserialization
 */

func TestSerializeWRQPacket(t *testing.T) {
	expected := []byte{0, 2, 'f', 'i', 'l', 'e', '.', 't', 'x', 't', 0, 'o', 'c', 't', 'e', 't', 0} // Opcode 2 (WRQ), Filename "file.txt", Mode "octet"

	inputMode, err := tftp.StringToMode("octet")
	if err != nil {
		t.Fatalf("Unexpected error converting mode string to Mode: %v", err)
	}

	inputFilename := "file.txt"

	connection := &tftp.TFTPConnection{
		Filename: &inputFilename,
		Mode:     &inputMode,
	}

	inputPacket := &tftp.TFTPPacket{
		Connection: *connection,
		Opcode:     2,
	}

	actual, err := inputPacket.Serialize()
	if err != nil {
		t.Fatalf("Unexpected error during serialization: %v", err)
	}
	if string(actual) != string(expected) {
		t.Fatalf("Expected bytes %v, got %v", expected, actual)
	}
}

func TestDeserializeWRQPacket(t *testing.T) {
	expectedFilename := "file.txt"
	expectedMode, err := tftp.StringToMode("octet")
	if err != nil {
		t.Fatalf("Unexpected error converting mode string to Mode: %v", err)
	}

	input := []byte{0, 2, 'f', 'i', 'l', 'e', '.', 't', 'x', 't', 0, 'o', 'c', 't', 'e', 't', 0}

	actualPacket, err := tftp.Deserialize(input)
	if err != nil {
		t.Fatalf("Unexpected error during deserialization: %v", err)
	}

	if actualPacket.Opcode != 2 {
		t.Fatalf("Expected Opcode 2, got %d", actualPacket.Opcode)
	}
	if *actualPacket.Connection.Filename != expectedFilename {
		t.Fatalf("Expected Filename '%s', got '%s'", expectedFilename, *actualPacket.Connection.Filename)
	}
	if *actualPacket.Connection.Mode != expectedMode {
		t.Fatalf("Expected Mode %d, got %d", expectedMode, *actualPacket.Connection.Mode)
	}
}

package tftp

import (
	"fmt"
	bytehelpers "networking/internal/byte_helpers"
	"networking/internal/logger"
)

func (t TFTPPacket) serializeInitialPacket() []byte {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))
	filename := []byte(*t.Connection.Filename)

	mode := *t.Connection.Mode
	modeBytes := bytehelpers.Uint16ToByteArray(uint16(mode))

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		filename,
		[]byte{0},
		modeBytes,
		[]byte{0},
	)

	return result
}

func (t TFTPPacket) serializeDataPacket() []byte {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))
	blockNumBytes := bytehelpers.Uint16ToByteArray(t.Connection.BlockCount)
	dataBytes := *t.Data

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		blockNumBytes,
		dataBytes,
	)

	return result
}

func (t TFTPPacket) serializeAckPacket() []byte {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))
	blockNumBytes := bytehelpers.Uint16ToByteArray(t.Connection.BlockCount)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		blockNumBytes,
	)

	return result
}

func (t TFTPPacket) serializeErrorPacket() []byte {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))
	errorCodeBytes := bytehelpers.Uint16ToByteArray(*t.ErrorCode)
	errorMsg, err := ErrorCodeToString(*t.ErrorCode)
	if err != nil {
		l := logger.GetLoggerFromContext(t.Connection.Ctx, &moduleName)
		errMsg := fmt.Sprintf("unknown error code: %d", *t.ErrorCode)
		l.Error(errMsg)
		return nil
	}

	errorMsgBytes := []byte(errorMsg)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		errorCodeBytes,
		errorMsgBytes,
		[]byte{0},
	)

	return result
}

func (t TFTPPacket) Serialize() []byte {
	switch t.Opcode {
	case 1:
		return t.serializeInitialPacket()
	case 2:
		return t.serializeInitialPacket()
	case 3:
		return t.serializeDataPacket()
	case 4:
		return t.serializeAckPacket()
	case 5:
		return t.serializeErrorPacket()
	}

	return nil
}

func deserializeOpcode(data []byte) (Opcode, error) {
	opcodeAsUint := bytehelpers.ByteArrayToUint16(data[0:2])
	opcode := Opcode(opcodeAsUint)
	_, err := OpcodeToString(opcode)
	if err != nil {
		err := fmt.Errorf("failed to deserialize initial packet: %v", err)
		return 0, err
	}

	return opcode, nil
}

func deserializeMode(startIndex int, data []byte) (Mode, error) {
	modeStr, _ := bytehelpers.GetNullTerminatedStringFromBytes(data[startIndex:])
	if modeStr == "" {
		return 0, fmt.Errorf("mode string is empty")
	}

	return StringToMode(modeStr)
}

func deserializeInitialPacket(data []byte) (*TFTPPacket, error) {
	opcode, err := deserializeOpcode(data)
	if err != nil {
		return nil, err
	}

	filename, endIndex := bytehelpers.GetNullTerminatedStringFromBytes(data[2:])
	mode, err := deserializeMode(endIndex, data)

	connection := &TFTPConnection{
		Filename: &filename,
		Mode:     &mode,
	}

	packet := &TFTPPacket{
		Connection: *connection,
		Opcode:     opcode,
	}

	return packet, nil
}

func deserializeDataPacket(data []byte) (*TFTPPacket, error) {
	opcode, err := deserializeOpcode(data)
	if err != nil {
		return nil, err
	}

	blockNum := bytehelpers.ByteArrayToUint16(data[2:4])
	dataBytes := data[4:]

	packet := &TFTPPacket{
		Opcode:   opcode,
		BlockNum: &blockNum,
		Data:     &dataBytes,
	}

	return packet, nil
}

func deserializeAckPacket(data []byte) (*TFTPPacket, error) {
	opcode, err := deserializeOpcode(data)
	if err != nil {
		return nil, err
	}

	blockNum := bytehelpers.ByteArrayToUint16(data[2:4])

	packet := &TFTPPacket{
		Opcode:   opcode,
		BlockNum: &blockNum,
	}

	return packet, nil
}

func deserializeErrorPacket(data []byte) (*TFTPPacket, error) {
	opcode, err := deserializeOpcode(data)
	if err != nil {
		return nil, err
	}

	errorCode := bytehelpers.ByteArrayToUint16(data[2:4])
	errorMsg, _ := bytehelpers.GetNullTerminatedStringFromBytes(data[4:])

	packet := &TFTPPacket{
		Opcode:    opcode,
		ErrorCode: &errorCode,
		ErrorMsg:  &errorMsg,
	}

	return packet, nil
}

func Deserialize(data []byte) (*TFTPPacket, error) {
	if len(data) < 2 {
		err := fmt.Errorf("data too short to contain opcode")
		return nil, err
	}

	opcode, err := deserializeOpcode(data)
	if err != nil {
		return nil, err
	}

	switch opcode {
	case 1, 2:
		return deserializeInitialPacket(data)
	case 3:
		return deserializeDataPacket(data)
	case 4:
		return deserializeAckPacket(data)
	case 5:
		return deserializeErrorPacket(data)
	}

	err = fmt.Errorf("unknown opcode: %d", opcode)
	return nil, err
}

package tftp

import (
	"fmt"
	bytehelpers "networking/internal/byte_helpers"
)

func (t TFTPPacket) checkFilename() error {
	if t.Connection.Filename == nil || *t.Connection.Filename == "" {
		return fmt.Errorf("filename is nil or empty")
	}

	return nil
}

func (t TFTPPacket) checkMode() error {
	if t.Connection.Mode == nil {
		return fmt.Errorf("mode is nil")
	}

	return nil
}

func (t TFTPPacket) checkBlockNum() error {
	if t.BlockNum == nil {
		return fmt.Errorf("cannot serialize DATA packet: Block number is nil")
	}

	return nil
}

func (t TFTPPacket) serializeInitialPacket() ([]byte, error) {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))

	if err := t.checkFilename(); err != nil {
		return nil, err
	}
	filename := []byte(*t.Connection.Filename)

	if err := t.checkMode(); err != nil {
		return nil, err
	}
	rawMode := *t.Connection.Mode
	modeString, err := ModeToString(uint8(rawMode))
	if err != nil {
		return nil, err
	}
	modeBytes := []byte(modeString)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		filename,
		[]byte{0},
		modeBytes,
		[]byte{0},
	)

	return result, nil
}

func (t TFTPPacket) serializeDataPacket() ([]byte, error) {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))

	if err := t.checkBlockNum(); err != nil {
		return nil, err
	}
	blockNumBytes := bytehelpers.Uint16ToByteArray(*t.BlockNum)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		blockNumBytes,
	)

	if t.Data != nil {
		result = bytehelpers.ConcatenateByteArrays(result, *t.Data)
	}

	return result, nil
}

func (t TFTPPacket) serializeAckPacket() ([]byte, error) {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))

	if t.BlockNum == nil {
		err := fmt.Errorf("cannot serialize ACK packet: Block number is nil")
		return nil, err
	}
	blockNumBytes := bytehelpers.Uint16ToByteArray(*t.BlockNum)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		blockNumBytes,
	)

	return result, nil
}

func (t TFTPPacket) serializeErrorPacket() ([]byte, error) {
	opcodeBytes := bytehelpers.Uint16ToByteArray(uint16(t.Opcode))

	if t.ErrorCode == nil {
		return nil, fmt.Errorf("cannot serialize ERROR packet: ErrorCode is nil")
	}
	errorCodeBytes := bytehelpers.Uint16ToByteArray(*t.ErrorCode)

	errorMsg, err := ErrorCodeToString(*t.ErrorCode)
	if err != nil {
		return nil, err
	}

	errorMsgBytes := []byte(errorMsg)

	result := bytehelpers.ConcatenateByteArrays(
		opcodeBytes,
		errorCodeBytes,
		errorMsgBytes,
		[]byte{0},
	)

	return result, nil
}

func (t TFTPPacket) Serialize() ([]byte, error) {
	if t.Opcode == 0 {
		return nil, fmt.Errorf("cannot serialize packet: Opcode is 0")
	}

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

	return nil, fmt.Errorf("cannot serialize packet: unknown Opcode %d", t.Opcode)
}

func deserializeOpcode(data []byte) (Opcode, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("data too short to contain opcode")
	}

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

	mode, err := deserializeMode(2+endIndex, data)
	if err != nil {
		return nil, err
	}

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

	if len(data) < 4 {
		return nil, fmt.Errorf("DATA packet too short to contain Block Number")
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

	if len(data) < 4 {
		return nil, fmt.Errorf("ACK packet too short to contain Block Number")
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

	if len(data) < 4 {
		return nil, fmt.Errorf("ERROR packet too short to contain Error Code")
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

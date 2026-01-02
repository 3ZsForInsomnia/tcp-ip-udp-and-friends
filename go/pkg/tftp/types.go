package tftp

import (
	"context"
	"fmt"
	"networking/internal/ports"
	"networking/internal/types"
)

var moduleName = "TFTP"

const maxDataLength = 512

const (
	RRQ  = 1
	WRQ  = 2
	DATA = 3
	ACK  = 4
	ERR  = 5
)

type Opcode uint16

func OpcodeToString(opcode Opcode) (string, error) {
	switch opcode {
	case 1:
		return "RRQ", nil
	case 2:
		return "WRQ", nil
	case 3:
		return "DATA", nil
	case 4:
		return "ACK", nil
	case 5:
		return "ERR", nil
	}

	err := fmt.Errorf("unknown opcode: %d", opcode)

	return "", err
}

func StringToOpcode(opcodeStr string) (Opcode, error) {
	switch opcodeStr {
	case "RRQ":
		return 1, nil
	case "WRQ":
		return 2, nil
	case "DATA":
		return 3, nil
	case "ACK":
		return 4, nil
	case "ERR":
		return 5, nil
	}

	err := fmt.Errorf("unknown opcode string: %s", opcodeStr)
	return 0, err
}

const (
	NETASCII = "netascii"
	OCTET    = "octet"
	MAIL     = "mail"
)

type Mode uint16

func ModeToString(mode uint8) (string, error) {
	switch mode {
	case 1:
		return NETASCII, nil
	case 2:
		return OCTET, nil
	case 3:
		return MAIL, nil
	}

	err := fmt.Errorf("unknown mode: %d", mode)
	return "", err
}

func StringToMode(modeStr string) (Mode, error) {
	switch modeStr {
	case NETASCII:
		return 1, nil
	case OCTET:
		return 2, nil
	case MAIL:
		return 3, nil
	}

	err := fmt.Errorf("unknown mode string: %s", modeStr)
	return 0, err
}

type TFTPConnection struct {
	Config     types.Config
	Connected  bool
	Filename   *string
	SourceTID  ports.Port
	DestTID    ports.Port
	Mode       *Mode
	Ctx        context.Context
	Data       *[]byte
	BlockCount uint16
}

const (
	NotDefined       = 0
	NotFound         = 1
	AccessViolation  = 2
	DiskFull         = 3
	IllegalOperation = 4
	UnknownTID       = 5
	FileExists       = 6
	NoSuchUser       = 7
)

func ErrorCodeToString(code uint16) (string, error) {
	switch code {
	case NotDefined:
		return "Not defined.", nil
	case NotFound:
		return "File not found.", nil
	case AccessViolation:
		return "Access violation.", nil
	case DiskFull:
		return "Disk full or allocation exceeded.", nil
	case IllegalOperation:
		return "Illegal TFTP operation.", nil
	case UnknownTID:
		return "Unknown transfer ID.", nil
	case FileExists:
		return "File already exists.", nil
	case NoSuchUser:
		return "No such user.", nil
	}

	return "", fmt.Errorf("unknown error code: %d", code)
}

type TFTPPacket struct {
	Connection TFTPConnection
	Opcode     Opcode
	Data       *[]byte
	ErrorCode  *uint16
	ErrorMsg   *string
	BlockNum   *uint16
}

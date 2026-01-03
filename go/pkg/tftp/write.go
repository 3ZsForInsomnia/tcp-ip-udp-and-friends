package tftp

import (
	"networking/internal/types"
	"networking/pkg/udp"
	"os"
)

func (t *TFTPConnection) OpenWriteConnection() error {
	writeRequest := TFTPPacket{
		Connection: *t,
		Opcode:     2, // WRQ
	}

	serializedRequest, err := writeRequest.Serialize()
	if err != nil {
		return err
	}

	err = udp.Send(types.Config(t.Config), serializedRequest)

	return err
}

func (t *TFTPConnection) SendFinalWriteAck() error {
	finalAck := TFTPPacket{
		Connection: *t,
		Opcode:     4, // ACK
	}

	serializedAck, err := finalAck.Serialize()
	if err != nil {
		return err
	}

	err = udp.Send(types.Config(t.Config), serializedAck)

	return err
}

func (t *TFTPConnection) CloseWriteConnection() error {
	err := t.SendFinalWriteAck()
	if err != nil {
		return err
	}

	t.Connected = false

	return nil
}

func (t *TFTPConnection) ReadSourceFile() ([]byte, error) {
	data, err := os.ReadFile(*t.Filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TFTPConnection) getNextPacketsBytes() []byte {
	currentBlockNum := t.BlockCount + 1

	startIndex := currentBlockNum * maxDataLength
	endIndex := (currentBlockNum + 1) * maxDataLength

	dataSlice := (*t.Data)[startIndex:endIndex]

	return dataSlice
}

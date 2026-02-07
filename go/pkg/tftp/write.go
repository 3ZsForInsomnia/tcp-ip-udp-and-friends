package tftp

import (
	"fmt"
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
	if err != nil {
		return err
	}

	t.Connected = true

	return nil
}

func (t *TFTPConnection) SendWritePacket(data []byte) error {
	packet := TFTPPacket{
		Connection: *t,
		Opcode:     3, // DATA
		Data:       &data,
		BlockNum:   &t.BlockCount,
	}

	serialized, err := packet.Serialize()
	if err != nil {
		return err
	}

	err = udp.Send(types.Config(t.Config), serialized)
	if err != nil {
		return err
	}

	t.BlockCount++

	return nil
}

func (t *TFTPConnection) CloseWriteConnection(data []byte) error {
	err := t.SendWritePacket(data)
	if err != nil {
		return err
	}

	t.Connected = false

	return nil
}

func (t *TFTPConnection) ReadSourceFile() error {
	data, err := os.ReadFile(*t.Filename)
	if err != nil {
		return err
	}

	t.Data = &data
	return nil
}

func (t *TFTPConnection) GetNextPacketsBytes() []byte {
	currentBlockNum := t.BlockCount + 1

	startIndex := currentBlockNum * maxDataLength
	endIndex := (currentBlockNum + 1) * maxDataLength

	dataSlice := (*t.Data)[startIndex:endIndex]

	return dataSlice
}

func (t *TFTPConnection) SendFile() error {
	err := t.ReadSourceFile()
	if err != nil {
		return err
	}
	t.BlockCount = 0

	c := make(chan udp.UDPGram)

	go udp.Listen(types.Config(t.Config), c)
	for gram := range c {
		packet, err := Deserialize(gram.Data)
		if err != nil {
			return err
		}

		opcode := packet.Opcode

		if opcode != 4 { // ACK
			err := fmt.Errorf("expected ACK packet, got opcode %d", opcode)
			return err
		}

		acknowledgedBlockNum := packet.BlockNum
		if acknowledgedBlockNum != &t.BlockCount {
			err := fmt.Errorf("expected ACK for block %d, got ACK for block %d", t.BlockCount, acknowledgedBlockNum)
			return err
		}

		dataSlice := t.GetNextPacketsBytes()
		if len(dataSlice) < maxDataLength {
			err := t.CloseWriteConnection(dataSlice)
			if err != nil {
				return err
			}
		}

		err = t.SendWritePacket(dataSlice)
		if err != nil {
			return err
		}
	}

	return nil
}

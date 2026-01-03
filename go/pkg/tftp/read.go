package tftp

import (
	"fmt"
	"networking/internal/ports"
	"networking/internal/types"
	"networking/pkg/udp"
	"os"
)

func (t *TFTPConnection) OpenReadConnection() error {
	randomPort := ports.GetRandomPort()
	randPortUint := uint16(randomPort)
	t.Config.SourcePort = &randPortUint

	firstPacket := TFTPPacket{
		Opcode:     1, // RRQ
		Connection: *t,
	}

	data, err := firstPacket.Serialize()
	if err != nil {
		return err
	}
	if data == nil {
		err := fmt.Errorf("failed to serialize packet for opening connection")
		return err
	}

	err = udp.Send(types.Config(t.Config), data)

	return err
}

func (t *TFTPConnection) sendFinalReadAck() error {
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

func (t *TFTPConnection) writeToFile(fileData []byte) error {
	filename := t.Filename
	err := os.WriteFile(*filename, fileData, 0644)

	return err
}

func (t *TFTPConnection) CloseReadConnection(data *[]byte) error {
	if data == nil {
		err := fmt.Errorf("no data to write to file on connection close")
		return err
	}

	err := t.sendFinalReadAck()
	if err != nil {
		return err
	}

	t.Connected = false

	err = t.writeToFile(*data)

	return err
}

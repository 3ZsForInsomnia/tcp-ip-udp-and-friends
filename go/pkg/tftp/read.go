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

	t.SourceTID = ports.Port(randomPort)
	t.DestTID = ports.Port(TftpDefaultPort)

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

	t.Connected = true

	return err
}

func (t *TFTPConnection) SendReadAck(blockNumber uint16) error {
	packet := TFTPPacket{
		Connection: *t,
		BlockNum:   &blockNumber,
		Opcode:     4, // ACK
	}

	serialized, err := packet.serializeAckPacket()
	if err != nil {
		return err
	}

	err = udp.Send(types.Config(t.Config), serialized)
	if err != nil {
		return err
	}

	return nil
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
	t.Connected = false

	return err
}

func (t *TFTPConnection) WriteToFile(fileData []byte) error {
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

	err = t.WriteToFile(*data)

	return err
}

func (t *TFTPConnection) ListenForReadData() error {
	c := make(chan udp.UDPGram)

	t.Config.SourcePort = (*uint16)(&t.SourceTID)
	go udp.Listen(types.Config(t.Config), c)

	data := make([]byte, 0)
	t.BlockCount = 1

	for gram := range c {
		packet, err := Deserialize(gram.Data)
		if err != nil {
			return err
		}

		t.DestTID = ports.Port(gram.SourcePort)

		if packet.Opcode == 3 { // DATA
			t.BlockCount++
			receivedData := *packet.Data
			data = append(data, receivedData...)

			if len(receivedData) < maxDataLength {
				t.Connected = false
				t.CloseReadConnection(&data)

				break
			} else {
				err := t.SendReadAck(*packet.BlockNum)
				if err != nil {
					return err
				}
			}

		} else if packet.Opcode == 5 { // ERROR
			errCode := packet.ErrorCode
			errMsg := *packet.ErrorMsg

			err := fmt.Errorf("received error from server: code %d, message %s", errCode, errMsg)

			return err

		} else {
			err := fmt.Errorf("unexpected opcode %d received", packet.Opcode)
			return err
		}
	}

	if len(data) == 0 {
		err := fmt.Errorf("no data received from server")
		return err
	}

	return nil
}

// Package tftp implements Trivial File Transfer Protocol (TFTP) functionalities.
package tftp

import (
	"context"
	"fmt"
	"networking/internal/logger"
	"networking/internal/ports"
	"networking/internal/types"
	"networking/pkg/udp"
	"os"
)

// TftpDefaultPort Selected this port to avoid conflicts with standard TFTP port 69
// The value is 2056 + 69, with 2056 being my "zero" port
const TftpDefaultPort = 2125

func StartTFTPServer(config types.Config, port ports.Port) error {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	portUint16 := uint16(port)
	config.ListenOnPort = &portUint16

	l.Info(fmt.Sprintf("Starting TFTP server on port %d", port))

	c := make(chan udp.UDPGram)
	go udp.Listen(config, c)

	for gram := range c {
		packet, err := Deserialize(gram.Data)
		if err != nil {
			l.Error(fmt.Sprintf("Failed to deserialize packet: %s", err.Error()))
			continue
		}

		clientPort := gram.SourcePort

		if packet.Opcode == RRQ {
			l.Info(fmt.Sprintf("Received RRQ for file '%s' from port %d", *packet.Connection.Filename, clientPort))
			err := HandleReadRequest(ctx, config, *packet.Connection.Filename, packet.Connection.Mode, ports.Port(clientPort))
			if err != nil {
				l.Error(fmt.Sprintf("Failed to handle RRQ: %s", err.Error()))
			} else {
				l.Info(fmt.Sprintf("Successfully sent file '%s' to client", *packet.Connection.Filename))
			}
		} else if packet.Opcode == WRQ {
			l.Info(fmt.Sprintf("Received WRQ for file '%s' from port %d", *packet.Connection.Filename, clientPort))
			err := HandleWriteRequest(ctx, config, *packet.Connection.Filename, packet.Connection.Mode, ports.Port(clientPort))
			if err != nil {
				l.Error(fmt.Sprintf("Failed to handle WRQ: %s", err.Error()))
			} else {
				l.Info(fmt.Sprintf("Successfully received file '%s' from client", *packet.Connection.Filename))
			}
		} else {
			l.Error(fmt.Sprintf("Received unexpected opcode %d as initial packet", packet.Opcode))
			continue
		}
	}

	return nil
}

func HandleReadRequest(ctx context.Context, config types.Config, filename string, mode *Mode, clientPort ports.Port) error {
	l := logger.GetLoggerFromContext(ctx)

	fileData, err := os.ReadFile(filename)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to read file '%s': %s", filename, err.Error()))
		return err
	}

	randomPort := ports.GetRandomPort()
	sourcePortUint16 := uint16(randomPort)
	clientPortUint16 := uint16(clientPort)
	config.SourcePort = &sourcePortUint16
	config.DestinationPort = &clientPortUint16

	blockNum := uint16(1)
	dataIndex := 0

	for {
		endIndex := dataIndex + maxDataLength
		if endIndex > len(fileData) {
			endIndex = len(fileData)
		}

		dataChunk := fileData[dataIndex:endIndex]
		
		packet := TFTPPacket{
			Opcode:   DATA,
			BlockNum: &blockNum,
			Data:     &dataChunk,
		}

		serialized, err := packet.Serialize()
		if err != nil {
			return err
		}

		err = udp.Send(config, serialized)
		if err != nil {
			return err
		}

		// Wait for ACK
		ackChan := make(chan udp.UDPGram)
		config.ListenOnPort = &sourcePortUint16
		go udp.Listen(config, ackChan)

		ackGram := <-ackChan
		ackPacket, err := Deserialize(ackGram.Data)
		if err != nil {
			return err
		}

		if ackPacket.Opcode != ACK {
			return fmt.Errorf("expected ACK, got opcode %d", ackPacket.Opcode)
		}

		if *ackPacket.BlockNum != blockNum {
			return fmt.Errorf("expected ACK for block %d, got ACK for block %d", blockNum, *ackPacket.BlockNum)
		}

		// Check if this was the last packet
		if len(dataChunk) < maxDataLength {
			break
		}

		blockNum++
		dataIndex = endIndex
	}

	return nil
}

func HandleWriteRequest(ctx context.Context, config types.Config, filename string, mode *Mode, clientPort ports.Port) error {
	l := logger.GetLoggerFromContext(ctx)

	randomPort := ports.GetRandomPort()
	sourcePortUint16 := uint16(randomPort)
	clientPortUint16 := uint16(clientPort)
	config.SourcePort = &sourcePortUint16
	config.DestinationPort = &clientPortUint16

	// Send initial ACK (block 0)
	blockNum := uint16(0)
	ackPacket := TFTPPacket{
		Opcode:   ACK,
		BlockNum: &blockNum,
	}

	serialized, err := ackPacket.Serialize()
	if err != nil {
		return err
	}

	err = udp.Send(config, serialized)
	if err != nil {
		return err
	}

	// Receive DATA packets
	dataChan := make(chan udp.UDPGram)
	config.ListenOnPort = &sourcePortUint16
	go udp.Listen(config, dataChan)

	var fileData []byte
	expectedBlockNum := uint16(1)

	for {
		dataGram := <-dataChan
		dataPacket, err := Deserialize(dataGram.Data)
		if err != nil {
			l.Error(fmt.Sprintf("Failed to deserialize DATA packet: %s", err.Error()))
			return err
		}

		if dataPacket.Opcode == ERR {
			return fmt.Errorf("received error from client: code %d, message %s", *dataPacket.ErrorCode, *dataPacket.ErrorMsg)
		}

		if dataPacket.Opcode != DATA {
			return fmt.Errorf("expected DATA packet, got opcode %d", dataPacket.Opcode)
		}

		if *dataPacket.BlockNum != expectedBlockNum {
			return fmt.Errorf("expected block %d, got block %d", expectedBlockNum, *dataPacket.BlockNum)
		}

		fileData = append(fileData, *dataPacket.Data...)

		// Send ACK
		ackPacket := TFTPPacket{
			Opcode:   ACK,
			BlockNum: dataPacket.BlockNum,
		}

		serialized, err := ackPacket.Serialize()
		if err != nil {
			return err
		}

		err = udp.Send(config, serialized)
		if err != nil {
			return err
		}

		// Check if this was the last packet
		if len(*dataPacket.Data) < maxDataLength {
			break
		}

		expectedBlockNum++
	}

	// Write file to disk
	err = os.WriteFile(filename, fileData, 0644)
	if err != nil {
		l.Error(fmt.Sprintf("Failed to write file '%s': %s", filename, err.Error()))
		return err
	}

	return nil
}

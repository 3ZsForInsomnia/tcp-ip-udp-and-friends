// Package udp provides functionality to create, send, and listen for UDP datagrams.
package udp

import (
	"context"
	"fmt"
	"net"

	"networking/internal/logger"
	"networking/internal/types"
)

var moduleName = "UDP"

func Send(config types.Config, data []byte) error {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	if config.DestinationPort == nil {
		destError := fmt.Errorf("DestinationPort is not set in config")
		l.Error(destError.Error())

		return destError
	}

	sourcePortStr := fmt.Sprint(*config.SourcePort)

	conn, err := net.ListenPacket("udp", "127.0.0.1:"+sourcePortStr)
	if err != nil {
		l.Error(err.Error())
		return err
	}
	defer conn.Close()

	datagram, err := NewUDPGram(&ctx, config.SourcePort, config.DestinationPort, &data)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	datagramBytes, err := datagram.CreateUDPGram(&ctx)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	addr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: int(*config.DestinationPort)}

	_, sendErr := conn.WriteTo(datagramBytes, &addr)
	if sendErr != nil {
		l.Error(sendErr.Error())
		return sendErr
	}

	return nil
}

func Listen(config types.Config, c chan UDPGram) {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	if config.ListenOnPort == nil {
		l.Error("ListenOnPort is not set in config")
		panic("ListenOnPort is not set in config")
	}

	destPortStr := fmt.Sprint(*config.ListenOnPort)

	conn, err := net.ListenPacket("udp", "127.0.0.1:"+destPortStr)
	if err != nil {
		l.Error(err.Error())
		panic(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			l.Error("Error closing connection: " + err.Error())
			panic(err)
		}
	}()

	buffer := make([]byte, 65535)
	for {
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			l.Error("Error reading incoming packet: " + err.Error())
		} else {
			// l.Info(fmt.Sprintf("Received %d bytes from %s", n, addr.String()))
			gram, err := ParseRawUDPGram(ctx, buffer[:n])
			if err != nil {
				l.Error("Error parsing UDP datagram: " + err.Error())
			}

			c <- *gram
		}
	}
}

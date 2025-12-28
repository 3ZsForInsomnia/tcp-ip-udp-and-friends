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

	conn, err := net.ListenPacket("ip4:17", "127.0.0.1")
	if err != nil {
		l.Error(err.Error())
		return err
	}
	defer conn.Close()

	datagram, err := NewUDPGram(&ctx, nil, config.DestinationPort, &data)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	datagramBytes, err := datagram.CreateUDPGram(&ctx)
	if err != nil {
		l.Error(err.Error())
		return err
	}

	addr := net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}

	_, sendErr := conn.WriteTo(datagramBytes, &addr)
	if sendErr != nil {
		l.Error(sendErr.Error())
		return sendErr
	}

	return nil
}

func Listen(config types.Config) error {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	if config.ListenOnPort == nil {
		defaultPort := uint16(8080)
		config.ListenOnPort = &defaultPort
		l.Info("ListenOnPort is not set in config, using default port 8080")
	}

	conn, err := net.ListenPacket("ip4:17", "127.0.0.1")
	if err != nil {
		l.Error(err.Error())
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 65535)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			l.Error("Error reading incoming packet: " + err.Error())
		} else {
			l.Info(fmt.Sprintf("Received %d bytes from %s", n, addr.String()))
			_, err := ParseRawUDPGram(ctx, buffer[:n])
			if err != nil {
				l.Error("Error parsing UDP datagram: " + err.Error())
			}
		}
	}
}

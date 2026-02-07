package main

import (
	"flag"
	"fmt"
	"networking/internal/logger"
	"networking/internal/ports"
	"networking/internal/types"
	"networking/pkg/udp"
)

var moduleName = "UDP"

func main() {
	var command string
	flag.StringVar(&command, "c", "send", "Command to execute: send or listen")
	if command != "send" && command != "listen" {
		panic("Invalid command. Use 'send' or 'listen', or 'help' for usage.")
	}

	var input string
	flag.StringVar(&input, "i", "Hello there, hi!", "Input string to send via UDP")
	if command == "send" && input == "" {
		panic("Input string cannot be empty when command is 'send'")
	}
	if command == "listen" {
		fmt.Println("Input string is not used when command is 'listen'. Did you mean to do this?")
	}

	var listenPort uint
	flag.UintVar(&listenPort, "l", 8080, "Port to listen on for UDP packets")
	if !ports.IsValidPort(uint16(listenPort)) {
		panic("Invalid listen port")
	}
	listenPortUint16 := uint16(listenPort)

	var destPort uint
	flag.UintVar(&destPort, "d", 8080, "Destination port to send UDP packets to")
	if !ports.IsValidPort(uint16(destPort)) {
		panic("Invalid destination port")
	}
	destPortUint16 := uint16(destPort)

	var sourcePort uint
	flag.UintVar(&sourcePort, "s", 8081, "Source port for UDP packets")
	sourcePortUint16 := uint16(sourcePort)

	var logLevel string
	flag.StringVar(&logLevel, "log", "INFO", "Logging level (INFO, WARN, ERROR)")
	level, err := logger.GetLogLevelFromString(logLevel)
	if err != nil {
		panic(err)
	}
	l := logger.NewLogger(&moduleName, level)

	flag.Parse()

	inputBytes := []byte(input)

	config := types.Config{
		ListenOnPort:    &listenPortUint16,
		DestinationPort: &destPortUint16,
		SourcePort:      &sourcePortUint16,
		LogLevel:        &level,
	}

	if command == "listen" {
		l.Info("Listening for UDP packets on port " + fmt.Sprint(listenPortUint16))
		c := make(chan udp.UDPGram)
		go udp.Listen(config, c)

		for gram := range c {
			l.Info("Received UDP packet:\n" + gram.String(2))
		}
	} else {
		l.Info("Sending UDP packet on port " + fmt.Sprint(destPortUint16))
		err := udp.Send(config, inputBytes)
		if err != nil {
			l.Error("Error sending UDP packet:" + err.Error())
		}
	}
}

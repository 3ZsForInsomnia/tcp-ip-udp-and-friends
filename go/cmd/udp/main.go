package main

import (
	"flag"
	"fmt"
	"networking/internal/logger"
	"networking/internal/ports"
	"networking/internal/types"
	"networking/pkg/udp"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "help", "Command to execute: send or listen")
	if command == "help" {
		flag.PrintDefaults()
	}
	if command != "send" && command != "listen" && command != "help" {
		panic("Invalid command. Use 'send' or 'listen', or 'help' for usage.")
	}

	var input string
	flag.StringVar(&input, "input", "Hello there, bitchfuck!", "Input string to send via UDP")
	if command == "send" && input == "" {
		panic("Input string cannot be empty when command is 'send'")
	}
	if command == "listen" {
		fmt.Println("Input string is not used when command is 'listen'")
	}

	var listenPort uint
	flag.UintVar(&listenPort, "listenPort", 8080, "Port to listen on for UDP packets")
	if !ports.IsValidPort(uint16(listenPort)) {
		panic("Invalid listen port")
	}
	listenPortUint16 := uint16(listenPort)

	var destPort uint
	flag.UintVar(&destPort, "destPort", 8081, "Destination port to send UDP packets to")
	if !ports.IsValidPort(uint16(destPort)) {
		panic("Invalid destination port")
	}
	destPortUint16 := uint16(destPort)

	var logLevel string
	flag.StringVar(&logLevel, "logLevel", "INFO", "Logging level (DEBUG, INFO, WARN, ERROR)")
	level, err := logger.GetLogLevelFromString(logLevel)
	if err != nil {
		panic(err)
	}

	inputBytes := []byte(input)

	config := types.Config{
		ListenOnPort:    &listenPortUint16,
		DestinationPort: &destPortUint16,
		LogLevel:        &level,
	}

	if command == "listen" {
		udp.Listen(config)
	} else {
		err := udp.Send(config, inputBytes)
		if err != nil {
			fmt.Println("Error sending UDP packet:", err)
		}
	}
}

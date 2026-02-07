package main

import (
	"flag"
	"fmt"
	"networking/internal/logger"
	"networking/internal/ports"
	"networking/internal/types"
	"networking/pkg/tftp"
)

var moduleName = "TFTP"

func main() {
	var command string
	flag.StringVar(&command, "c", "server", "Command to execute: server, get, or write")
	if command != "server" && command != "get" && command != "write" {
		panic("Invalid command. Use 'server', 'get', or 'write'")
	}

	var filename string
	flag.StringVar(&filename, "f", "", "Filename for get/write operations")
	if (command == "get" || command == "write") && filename == "" {
		panic("Filename is required for get/write commands")
	}

	var modeStr string
	flag.StringVar(&modeStr, "m", "octet", "Transfer mode: netascii, octet, or mail")

	var destPort uint
	flag.UintVar(&destPort, "d", tftp.TftpDefaultPort, "Destination port for get/write operations")
	if !ports.IsValidPort(uint16(destPort)) {
		panic("Invalid destination port")
	}
	destPortUint16 := uint16(destPort)

	var sourcePort uint
	flag.UintVar(&sourcePort, "s", 0, "Source port (0 for random)")
	var sourcePortUint16 uint16
	if sourcePort == 0 {
		sourcePortUint16 = uint16(ports.GetRandomPort())
	} else {
		if !ports.IsValidPort(uint16(sourcePort)) {
			panic("Invalid source port")
		}
		sourcePortUint16 = uint16(sourcePort)
	}

	var logLevel string
	flag.StringVar(&logLevel, "log", "INFO", "Logging level (INFO, WARN, ERROR)")
	level, err := logger.GetLogLevelFromString(logLevel)
	if err != nil {
		panic(err)
	}
	l := logger.NewLogger(&moduleName, level)

	flag.Parse()

	config := types.Config{
		DestinationPort: &destPortUint16,
		SourcePort:      &sourcePortUint16,
		LogLevel:        &level,
	}

	if command == "server" {
		l.Info(fmt.Sprintf("Starting TFTP server on port %d", tftp.TftpDefaultPort))
		err := tftp.StartTFTPServer(config, ports.Port(tftp.TftpDefaultPort))
		if err != nil {
			l.Error("Server error: " + err.Error())
			panic(err)
		}
	} else if command == "get" {
		l.Info(fmt.Sprintf("Getting file '%s' from port %d", filename, destPort))
		
		mode, err := tftp.StringToMode(modeStr)
		if err != nil {
			l.Error("Invalid mode: " + err.Error())
			panic(err)
		}

		err = tftp.GetFile(config, filename, &mode)
		if err != nil {
			l.Error("Error getting file: " + err.Error())
			panic(err)
		}
		l.Info(fmt.Sprintf("Successfully retrieved file '%s'", filename))
	} else if command == "write" {
		l.Info(fmt.Sprintf("Writing file '%s' to port %d", filename, destPort))
		
		mode, err := tftp.StringToMode(modeStr)
		if err != nil {
			l.Error("Invalid mode: " + err.Error())
			panic(err)
		}

		err = tftp.WriteFile(config, filename, &mode)
		if err != nil {
			l.Error("Error writing file: " + err.Error())
			panic(err)
		}
		l.Info(fmt.Sprintf("Successfully sent file '%s'", filename))
	}
}


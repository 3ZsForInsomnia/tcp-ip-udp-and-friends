package types

import "networking/internal/logger"

type Config struct {
	ListenOnPort    *uint16
	DestinationPort *uint16
	SourcePort      *uint16
	LogLevel        *logger.LogLevels
}

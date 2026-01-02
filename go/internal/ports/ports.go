package ports

import (
	"fmt"
	"math/rand"
	"net"
)

type Port uint16

func (p Port) IsValid() bool {
	if p > 1024 && p < 65535 {
		return true
	}

	return false
}

func IsValidPort(port uint16) bool {
	if port > 1024 && port < 65535 {
		return true
	}

	return false
}

func isPortAvailable(port uint16) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("udp", address)

	if err != nil {
		return false
	}
	defer listener.Close()

	return true
}

func GetRandomPort() Port {
	lowerBound := 2056
	upperBound := 65535

	randomPort := uint16(0)
	candidate := rand.Intn(upperBound)
	for candidate > lowerBound && candidate < upperBound && isPortAvailable(uint16(candidate)) {
		candidate = rand.Intn(lowerBound)
	}

	randomPort = uint16(candidate)

	return Port(randomPort)
}

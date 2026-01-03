// Package tftp implements Trivial File Transfer Protocol (TFTP) functionalities.
package tftp

import (
	"networking/internal/ports"
	"networking/internal/types"
)

// Server to listen goes here
// Allow sending read/write requests to TFTP servers

func StartTFTPServer(config types.Config, port ports.Port) error {
	return nil
}

func ReadFileFromTFTPServer(config types.Config, filename string) ([]byte, error) {
	return nil, nil
}

func WriteFileToTFTPServer(config types.Config, filename string, data []byte) error {
	return nil
}

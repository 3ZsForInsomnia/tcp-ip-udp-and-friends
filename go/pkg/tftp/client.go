// Package tftp implements Trivial File Transfer Protocol (TFTP) functionalities.
package tftp

import (
	"context"
	"networking/internal/logger"
	"networking/internal/ports"
	"networking/internal/types"
)

func GetFile(config types.Config, filename string, mode *Mode) error {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	var modeVal Mode
	if mode == nil {
		modeVal = GetDefaultMode()
	}

	data := make([]byte, 0)
	t := TFTPConnection{
		Config:    config,
		Ctx:       ctx,
		Filename:  &filename,
		Connected: false,
		DestTID:   ports.Port(TftpDefaultPort),
		Mode:      &modeVal,
		Data:      &data,
	}

	err := t.OpenReadConnection()
	if err != nil {
		l.Error(err.Error())
		return err
	}

	err = t.ListenForReadData()
	if err != nil {
		l.Error(err.Error())
		return err
	}

	return nil
}

func WriteFile(config types.Config, filename string, mode *Mode) error {
	l := logger.NewLogger(&moduleName, *config.LogLevel)
	ctx := l.WithLogger(context.Background())

	var modeVal Mode
	if mode == nil {
		modeVal = GetDefaultMode()
	}

	t := TFTPConnection{
		Config:    config,
		Ctx:       ctx,
		Filename:  &filename,
		Connected: false,
		DestTID:   ports.Port(TftpDefaultPort),
		Mode:      &modeVal,
		Data:      nil,
	}

	err := t.OpenWriteConnection()
	if err != nil {
		l.Error(err.Error())
		return err
	}

	err = t.SendFile()
	if err != nil {
		l.Error(err.Error())
		return err
	}

	return nil
}

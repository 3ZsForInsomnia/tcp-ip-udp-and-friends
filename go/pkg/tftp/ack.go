package tftp

import "fmt"

func (t *TFTPConnection) CreateAck(packet *TFTPPacket) (*TFTPPacket, error) {
	if !t.Connected {
		if *t.Config.SourcePort == 0 {
			err := fmt.Errorf("TFTP connection not established for read")
			return nil, err
		} else {
			portString := fmt.Sprint(*t.Config.SourcePort)
			err := fmt.Errorf("TFTP connection not established for read; using source port: %s", portString)
			return nil, err
		}
	}

	if packet != nil && packet.Data != nil {
		readAck := TFTPPacket{
			Connection: *t,
			Opcode:     4, // ACK
		}

		if len(*packet.Data) >= maxDataLength {
			err := t.closeReadConnection(t.Data)
			if err != nil {
				return nil, err
			}
		}

		return &readAck, nil
	}

	err := fmt.Errorf("cannot create read ack packet from nil packet")

	return nil, err
}

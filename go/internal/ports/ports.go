package ports

func IsValidPort(port uint16) bool {
	if port > 1024 && port < 65535 {
		return true
	}

	return false
}

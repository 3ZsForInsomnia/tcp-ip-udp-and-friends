package bytehelpers

func Uint16ToByteArray(value uint16) []byte {
	highByte := byte(value >> 8)
	lowByte := byte(value & 0xff)

	return []byte{highByte, lowByte}
}

func ByteArrayToUint16(data []byte) uint16 {
	highByte := uint16(data[0]) << 8
	lowByte := uint16(data[1])

	return highByte | lowByte
}

func ConcatenateByteArrays(byteArrays ...[]byte) []byte {
	result := []byte{}

	for _, byteArray := range byteArrays {
		result = append(result, byteArray...)
	}

	return result
}

func AreByteArraysEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

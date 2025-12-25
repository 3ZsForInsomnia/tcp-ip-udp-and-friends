package bytehelpers

const maxUnsignedShort = 0xFFFF

func CreateOnesComplementChecksum(data []byte) uint16 {
	if len(data) == 0 {
		return maxUnsignedShort
	}

	if len(data) == 1 {
		return ^uint16(data[0])
	}

	sum := getSumOfData(data)
	complement := ^(sum)

	return complement
}

func getSumOfData(data []byte) uint16 {
	// We make this uint32 to avoid overflow when adding carried over bits
	sum := uint32(0)

	for i := 0; i < len(data)-1; i += 2 {
		first := uint16(data[i]) << 8
		second := uint16(data[i+1])
		currentWord := first | second

		carriedSum := carryAroundAdd(uint16(sum), currentWord)

		sum = uint32(carriedSum)
	}

	if len(data)%2 == 1 {
		end := len(data) - 1
		currentWord := uint16(data[end])

		// The odd byte is treated as the high byte of a 16-bit word, with the low byte being zero
		carriedSum := carryAroundAdd(uint16(sum), currentWord<<8)

		sum = uint32(carriedSum)
	}

	return uint16(sum)
}

func carryAroundAdd(a, b uint16) uint16 {
	result := uint32(a) + uint32(b)

	if result > maxUnsignedShort {
		// This gets us the last bit_s_ that overflowed - we re removing 16 bits out of the _32_ bits of `result`
		carry := result >> 16

		// This is required because we convert to uint32 above! It makes it so we ignore anything above 16 bits
		result = (result & maxUnsignedShort)

		result += carry
	}

	return uint16(result)
}

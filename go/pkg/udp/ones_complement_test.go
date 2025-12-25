package udp

import (
	"testing"
)

func Test_CreateOnesComplementChecksum_HappyPath(t *testing.T) {
	expected := uint16(0xD5D2)
	input := []byte{97, 98, 99, 100, 101, 102} // "abcdef"

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_CreateOnesComplementChecksum_EmptyInput(t *testing.T) {
	expected := uint16(65535)
	input := []byte{}

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_CreateOnesComplementChecksum_SingleByteInput(t *testing.T) {
	expected := uint16(65280)
	input := []byte{255}

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_CreateOnesComplementChecksum_OddLengthInput(t *testing.T) {
	expected := uint16(0x6ED2)
	input := []byte{97, 98, 99, 100, 101, 102, 103} // Odd length input, "abcdefg"

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_CreateOnesComplementChecksum_AllZeroBytes(t *testing.T) {
	expected := uint16(65535)
	input := []byte{0, 0, 0, 0, 0, 0} // All zero bytes

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_CreateOnesComplementChecksum_MaxByteValues(t *testing.T) {
	expected := uint16(0)
	input := []byte{255, 255, 255, 255} // All bytes are 255

	actual := CreateOnesComplementChecksum(input)

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

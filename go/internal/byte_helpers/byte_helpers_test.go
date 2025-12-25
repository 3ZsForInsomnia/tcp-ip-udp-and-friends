package bytehelpers

import (
	"bytes"
	"testing"
)

func Test_Uint16ToByteArray_HappyPath(t *testing.T) {
	expected := []byte{0x1A, 0x2B}

	input := uint16(0x1A2B)
	actual := Uint16ToByteArray(input)

	if !bytes.Equal(actual, expected) {
		t.Errorf("Uint16ToByteArray(%d) = %v; want %v", input, actual, expected)
	}
}

func Test_Uint16ToByteArray_ZeroValue(t *testing.T) {
	expected := []byte{0x00, 0x00}

	input := uint16(0)
	actual := Uint16ToByteArray(input)

	if !bytes.Equal(actual, expected) {
		t.Errorf("Uint16ToByteArray(%d) = %v; want %v", input, actual, expected)
	}
}

func Test_ByteArrayToUint16_HappyPath(t *testing.T) {
	expected := uint16(0x1A2B)

	input := []byte{0x1A, 0x2B}
	actual := ByteArrayToUint16(input)

	if actual != expected {
		t.Errorf("ByteArrayToUint16(%v) = %d; want %d", input, actual, expected)
	}
}

func Test_ByteArrayToUint16_ZeroValue(t *testing.T) {
	expected := uint16(0)

	input := []byte{0x00, 0x00}
	actual := ByteArrayToUint16(input)

	if actual != expected {
		t.Errorf("ByteArrayToUint16(%v) = %d; want %d", input, actual, expected)
	}
}

func Test_ConcatenateByteArrays_HappyPath(t *testing.T) {
	expected := []byte{0x1A, 0x2B, 0x3C, 0x4D}

	input1 := []byte{0x1A, 0x2B}
	actual := ConcatenateByteArrays(input1, []byte{0x3C, 0x4D})

	if !bytes.Equal(actual, expected) {
		t.Errorf("ConcatenateByteArrays(...) = %v; want %v", actual, expected)
	}
}

func Test_ConcatenateByteArrays_EmptyArrays(t *testing.T) {
	expected := []byte{}

	input1 := []byte{}
	actual := ConcatenateByteArrays(input1, []byte{})

	if !bytes.Equal(actual, expected) {
		t.Errorf("ConcatenateByteArrays(...) = %v; want %v", actual, expected)
	}
}

func Test_AreByteArraysEqual_EqualArrays(t *testing.T) {
	expected := true

	input1 := []byte{0x1A, 0x2B, 0x3C}
	input2 := []byte{0x1A, 0x2B, 0x3C}
	actual := AreByteArraysEqual(input1, input2)

	if actual != expected {
		t.Errorf("AreByteArraysEqual(%v, %v) = %v; want %v", input1, input2, actual, expected)
	}
}

func Test_AreByteArraysEqual_UnequalArrays(t *testing.T) {
	expected := false

	input1 := []byte{0x1A, 0x2B, 0x3C}
	input2 := []byte{0x1A, 0x2B, 0x4D}
	actual := AreByteArraysEqual(input1, input2)

	if actual != expected {
		t.Errorf("AreByteArraysEqual(%v, %v) = %v; want %v", input1, input2, actual, expected)
	}
}

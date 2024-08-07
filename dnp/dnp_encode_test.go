package dnp

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type element struct {
	plain   string
	encoded string
}

var values = [15]element{
	{"MyModel", "MyModel"},
	{"MySerial", "MySerial"},
	{"MyProductCode", "MyProductCode"},
	{"ABC DEF", "ABC,20DEF"},
	{"123.456/3&8", "123.456,2F3,268"},
	{"123456#33", "123456,2333"},
	{"ABC@home", "ABC,40home"},
	{"ABC*33<4", "ABC,2A33,3C4"},
	{"20123.4", "20123.4"},
	{"a/asd asd/dddd", "a,2Fasd,20asd,2Fdddd"},
	{"¥", ",C2,A5"},
	{"a¥,bא, ٺ,cD", "a,C2,A5,2Cb,D7,90,2C,20,D9,BA,2CcD"},
	{"🂢🂣🂮🂹🃱", ",F0,9F,82,A2,F0,9F,82,A3,F0,9F,82,AE,F0,9F,82,B9,F0,9F,83,B1"},
	{"🫒🪕", ",F0,9F,AB,92,F0,9F,AA,95"},
	{"aæc", "a,C3,A6c"},
}

func TestDNPEncode(t *testing.T) {
	for _, value := range values {
		encoded := Encode(value.plain)
		assert.Equal(t, value.encoded, encoded, "Failed to encode value of: %s", value.plain)
	}
}

func TestDNPDecode(t *testing.T) {
	assertDecode := func(value element) {
		decoded, err := Decode(value.encoded)
		require.NoError(t, err)
		assert.Equal(t, value.plain, decoded, "Failed to decode value of: %s", value.encoded)
	}

	for _, value := range values {
		assertDecode(value)
	}

	assertDecode(element{"123.456/3&8", "123.456,2f3,268"})
}

func TestDNPDecodeError(t *testing.T) {
	decoded, err := Decode(",,")
	require.Error(t, err)
	assert.Equal(t, "invalid DNP escape \",,\"", err.Error())
	assert.Equal(t, "", decoded)
}

func TestIsHexWithSpecialCharacterInput(t *testing.T) {
	noHex := ishex('*')
	assert.False(t, noHex)

	unhex := unhex('*')
	assert.Equal(t, byte(0), unhex)
}

func TestDNPEncodeWithNoEscapingRequired(t *testing.T) {
	input := "MyModel"
	expected := "MyModel"
	encoded := Encode(input)
	assert.Equal(t, expected, encoded, "Failed to encode value of: %s", input)
}

func TestDNPEncodeWithEscapingRequired(t *testing.T) {
	input := "ABC DEF"
	expected := "ABC,20DEF"
	encoded := Encode(input)
	assert.Equal(t, expected, encoded, "Failed to encode value of: %s", input)
}

func TestDNPDecodeWithNoEscapingRequired(t *testing.T) {
	input := "MyModel"
	expected := "MyModel"
	decoded, err := Decode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, decoded, "Failed to decode value of: %s", input)
}

func TestDNPDecodeWithEscapingRequired(t *testing.T) {
	input := "ABC,20DEF"
	expected := "ABC DEF"
	decoded, err := Decode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, decoded, "Failed to decode value of: %s", input)
}

func TestDNPDecodeWithInvalidEscapeSequence(t *testing.T) {
	input := "ABC,2GDEF"
	_, err := Decode(input)
	require.Error(t, err)
}

func TestDNPEncodeWithSpecialCharacters(t *testing.T) {
	input := "aæc"
	expected := "a,C3,A6c"
	encoded := Encode(input)
	assert.Equal(t, expected, encoded, "Failed to encode value of: %s", input)
}

func TestDNPDecodeWithSpecialCharacters(t *testing.T) {
	input := "a,C3,A6c"
	expected := "aæc"
	decoded, err := Decode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, decoded, "Failed to decode value of: %s", input)
}

func TestDNPEncodeLongSequence(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789,"
	expected := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789,2C"
	encoded := Encode(input)
	assert.Equal(t, expected, encoded, "Failed to decode value of: %s", input)
}

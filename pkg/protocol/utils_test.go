package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SliceEqual(t *testing.T) {
	bytes := []byte{0x53, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x54, 0x65, 0x78, 0x74, 0x31, 0x32, 0x33, 0x34}

	assert.True(t, SliceEqual(bytes, "SampleText"))
}

func Test_SliceEqual_NotEqual(t *testing.T) {
	bytes := []byte{0x53, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x54, 0x65, 0x78, 0x74, 0x31, 0x32, 0x33, 0x34}
	assert.False(t, SliceEqual(bytes, "SampleTextNot"))
}

func Test_SliceEqual_Short(t *testing.T) {
	bytes := []byte{0x53, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x54, 0x65, 0x78}
	assert.False(t, SliceEqual(bytes, "SampleText"))
}

func Test_TokenizeString_Simple(t *testing.T) {
	input := "test test123 test3"
	output := TokenizeString([]byte(input))

	assert.Equal(t, 3, len(output))
	assert.Equal(t, "test", output[0])
	assert.Equal(t, "test123", output[1])
	assert.Equal(t, "test3", output[2])
}

func Test_TokenizeString_Quotes(t *testing.T) {
	input := "test \"test123 test3\""
	output := TokenizeString([]byte(input))

	assert.Equal(t, 2, len(output))
	assert.Equal(t, "test", output[0])
	assert.Equal(t, "test123 test3", output[1])
}

func Test_TokenizeString_MultipleQuotes(t *testing.T) {
	input := "test \"test123 test3\" \"test4\""
	output := TokenizeString([]byte(input))

	assert.Equal(t, 3, len(output))
	assert.Equal(t, "test", output[0])
	assert.Equal(t, "test123 test3", output[1])
	assert.Equal(t, "test4", output[2])
}

func Test_TokenizeString_BlockComment(t *testing.T) {
	input := "test /* Not Here */ \"test123 test3\" \"test4\""
	output := TokenizeString([]byte(input))

	assert.Equal(t, 3, len(output))
	assert.Equal(t, "test", output[0])
	assert.Equal(t, "test123 test3", output[1])
	assert.Equal(t, "test4", output[2])
}

func Test_TokenizeString_Comment(t *testing.T) {
	input := "test \"test123 test3\" \"test4\" //Testing"
	output := TokenizeString([]byte(input))

	assert.Equal(t, 3, len(output))
	assert.Equal(t, "test", output[0])
	assert.Equal(t, "test123 test3", output[1])
	assert.Equal(t, "test4", output[2])
}

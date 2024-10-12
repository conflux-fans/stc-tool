package zkclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolsToUint64(t *testing.T) {
	bools := []bool{true, false, true, true}
	assert.Equal(t, uint64(0b1011), BoolsToUint64(bools))
}

func TestReverseBools(t *testing.T) {
	bools := []bool{true, false, true, true}
	assert.Equal(t, []bool{true, true, false, true}, ReverseBools(bools))
}

func TestGenVcInputPath(t *testing.T) {
	bools := []bool{true, false, true, true}
	assert.Equal(t, uint64(0b0010), BoolsToUint64(ReverseBools(InvertBools(bools))))

	bools = []bool{
		false,
		true,
		true,
		true,
		false,
		false,
		true,
		true,
		false,
		true,
	}
	assert.Equal(t, uint64(0b0010), BoolsToUint64(ReverseBools(InvertBools(bools))))
	assert.Equal(t, uint64(0b0010), BoolsToUint64((InvertBools(bools))))
	assert.Equal(t, uint64(0b0010), BoolsToUint64(ReverseBools(InvertBools(bools))))
}

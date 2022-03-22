package rust

import (
	"testing"

	"gotest.tools/assert"
)

func TestIsBaseType(t *testing.T) {
	r := IsBaseType("H256")
	assert.Equal(t, r, true)
}

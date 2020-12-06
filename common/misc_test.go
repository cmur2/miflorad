package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	assert.Equal(t, "", MifloraGetAlphaNumericID(""))
	assert.Equal(t, "1234567890ab", MifloraGetAlphaNumericID("12:34:56:78:90:ab"))
}

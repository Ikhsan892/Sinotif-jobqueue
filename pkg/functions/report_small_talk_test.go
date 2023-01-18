package functions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const samplePayload string = "{ \"user_id\" : 1, \"student_id\" : 2 }"

func TestParsingPayload(t *testing.T) {
	data, err := ParsePayload(samplePayload)

	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, uint(1), data.UserId)
}

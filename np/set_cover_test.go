package np

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCoverSmall(t *testing.T) {
	assert := assert.New(t)
	sets := [][]string{
		{"a", "b", "c", "d", "e", "f"},
		{"e", "f", "h", "i"},
		{"a", "d", "g", "j"},
		{"b", "e", "h", "k"},
		{"c", "f", "i", "l"},
		{"j", "k", "l"},
	}
	setCover := MinimumSetCover(sets)
	assert.Equal(3, len(setCover))
}

func TestSetCornerCase(t *testing.T) {
	assert := assert.New(t)
	var sets [][]string
	setCover := MinimumSetCover(sets)
	assert.Equal(0, len(setCover))

	sets = append(sets, []string{})
	setCover = MinimumSetCover(sets)
	assert.Equal(0, len(setCover))

	sets = append(sets, []string{"a"})
	sets = append(sets, []string{"a"})
	sets = append(sets, []string{"a"})
	setCover = MinimumSetCover(sets)
	assert.Equal(1, len(setCover))

	sets = append(sets, []string{"b"})
	setCover = MinimumSetCover(sets)
	assert.Equal(2, len(setCover))

	sets = append(sets, []string{"a", "b"})
	setCover = MinimumSetCover(sets)
	assert.Equal(1, len(setCover))
}

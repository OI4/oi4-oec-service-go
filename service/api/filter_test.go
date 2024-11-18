package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFilter(t *testing.T) {
	var filter Filter
	filter = NewStringFilter("test")
	assert.Equal(t, "test", filter.String())
}

package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFilter(t *testing.T) {
	var filterA Filter = "test"
	var filterB Filter = "test"
	assert.True(t, FilterEquals(&filterA, &filterB))
}

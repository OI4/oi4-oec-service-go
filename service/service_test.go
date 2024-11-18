package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService(t *testing.T) {
	text := GetLocalizedText()
	fmt.Println(text)
	assert.NotNil(t, text)
}

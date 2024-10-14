package api

import (
	"testing"
)

func TestParseMethodType_ValidMethodTypes(t *testing.T) {
	validMethodTypes := []string{"Get", "Pub", "Set", "Del", "Call", "Reply"}

	for _, methodType := range validMethodTypes {
		_, err := ParseMethodType(methodType)
		if err != nil {
			t.Errorf("Expected no error for valid method type %s, got %v", methodType, err)
		}
	}
}

func TestParseMethodType_InvalidMethodType(t *testing.T) {
	_, err := ParseMethodType("InvalidMethodType")
	if err == nil {
		t.Errorf("Expected error for invalid method type, got nil")
	}
}

func TestParseMethodType_EmptyString(t *testing.T) {
	_, err := ParseMethodType("")
	if err == nil {
		t.Errorf("Expected error for empty string, got nil")
	}
}

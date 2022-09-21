package utils

import "testing"

func TestDNPEncode(t *testing.T) {
	if DNPEncode("123456#33") != "123456,2333" {
		t.Error("invalid encoding")
	}
	if DNPEncode("123.456/3&8") != "123.456,2F3,268" {
		t.Error("invalid encoding")
	}
}

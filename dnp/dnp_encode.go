// Package dnp provides encoding and decoding functions for DNP strings.
// The DNP is basically a URL encoding, replacing the '%' with an ',',
// so that the resulting string is  DIN SPEC 91406 compliant
//
// The basic encoding and decoding is taken from the net/url package.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package dnp

import (
	"strconv"
	"strings"
)

const escapeChar = ','

const upperHex = "0123456789ABCDEF"

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

type EscapeError string

func (e EscapeError) Error() string {
	return "invalid DNP escape " + strconv.Quote(string(e))
}

// Return true if the specified character should be escaped when
// appearing in a DNP string, according to DIN SPEC 91406.
func shouldEscape(c byte) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@', ' ': // §2.2 Reserved characters (reserved)
		return true
	}

	// Everything else must be escaped.
	return true
}

// unescape unescapes a DNP string
func unescape(s string) (string, error) {
	// Count ',', check that they're well-formed.
	n := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case escapeChar:
			n++
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[:3]
				}
				return "", EscapeError(s)
			}
			i += 3
		default:
			i++
		}
	}

	if n == 0 {
		return s, nil
	}

	var t strings.Builder
	t.Grow(len(s) - 2*n)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case escapeChar:
			t.WriteByte(unhex(s[i+1])<<4 | unhex(s[i+2]))
			i += 2
		default:
			t.WriteByte(s[i])
		}
	}
	return t.String(), nil
}

func escape(s string) string {
	hexCount := 0
	//spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = escapeChar
			t[j+1] = upperHex[c>>4]
			t[j+2] = upperHex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

// Encode escapes the string according to the DNP encoding.
// The DNP is basically a URL encoding, replacing the '%' with an ',',
// so that the resulting string is  DIN SPEC 91406 compliant
func Encode(input string) string {
	return escape(input)
}

// Decode unescapes a DNP encoded string.
// The DNP is basically a URL encoding, replacing the '%' with an ',',
// so that the resulting string is  DIN SPEC 91406 compliant
func Decode(input string) (string, error) {
	return unescape(input)
}

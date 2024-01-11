package validators

import (
	"net"
	"regexp"
)

// Validates a string is a valid email address
// This regular expression comes from the HTML5 specs and is used by web browsers to validate "email" input fields
var emailRe = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_\x60{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// Email returns true if the string is a valid email address.
func Email(val string) bool {
	return emailRe.MatchString(val)
}

// ItemID returns true if a string is a valid "NanoID" as generated by the "go-nanoid" library and used in this app.
func ItemID(val string) bool {
	// Validate the format is base64-url, 21-characters
	return Base64URL(val, 21)
}

// Base64URL returns true if the argument is a string that matches [a-zA-Z0-9-_] and is of the given length
func Base64URL(val string, expectLen int) bool {
	if len(val) != expectLen {
		return false
	}

	for i := 0; i < expectLen; i++ {
		valid := (val[i] >= 'a' && val[i] <= 'z') ||
			(val[i] >= 'A' && val[i] <= 'Z') ||
			(val[i] >= '0' && val[i] <= '9') ||
			val[i] == '_' || val[i] == '-'
		if !valid {
			return false
		}
	}

	return true
}

// IsIP returns true if the argument is a valid IPv4 or IPv6
func IsIP(val string) bool {
	return net.ParseIP(val) != nil
}

// IsHostname returns true if a string is a valid hostname.
func IsHostname(s string) bool {
	// Source: https://github.com/golang/go/blob/go1.21.6/src/net/dnsclient.go
	// Copyright 2009 The Go Authors
	// License: BSD (https://github.com/golang/go/blob/go1.21.6/LICENSE)

	// The root domain name is valid. See golang.org/issue/45715.
	if s == "." {
		return true
	}

	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return nonNumeric
}

package utils

import (
	"strconv"

	"github.com/martinlindhe/base36"
	//"strings"
)

// Base36Encoder returns base 36 of n.
func Base36Encoder(n int64) string {
	shortURL := strconv.FormatInt(n, 36)
	return shortURL
}

// Base36Decoder returns base 10 s (s is based on 36).
func Base36Decoder(s string) uint64 {
	id := base36.Decode(s)
	return id
}

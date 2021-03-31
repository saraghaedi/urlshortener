package functions

import (
	"github.com/martinlindhe/base36"
	"strconv"
	//"strings"
)
// afterDomain := strings.SplitN(bigurl, "/", 1)

// encode function is used to convert id to shorted URL.
func URLEncoder(i int64) string  {
	shortURL := strconv.FormatInt(i, 36)
	return shortURL
}

// decode function is used to convert shorted URL to id.
func URLDecoder(s string) uint64  {
	id := base36.Decode(s)
	return id
}





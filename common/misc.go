package common

import (
	"regexp"
	"strings"
)

const peripheralAddressAllowedChars = "[^a-z0-9]+"

var peripheralAddressAllowedCharsPattern = regexp.MustCompile(peripheralAddressAllowedChars)

func MifloraGetAlphaNumericID(peripheralAddress string) string {
	return peripheralAddressAllowedCharsPattern.ReplaceAllString(strings.ToLower(peripheralAddress), "")
}

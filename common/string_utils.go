package common

import (
	"bytes"
	"strings"
)

func MakeRange(min, max int) []int {
	num := make([]int, max-min+1)
	for i := range num {
		num[i] = min + i
	}
	return num
}

func isInRange(r []int, item int) bool {
	for i := range r {
		if i == item {
			return true
		}
	}
	return false
}

func ContainsInStringArray(in []string, item string) bool {
	for _, a := range in {
		if a == item {
			return true
		}
	}
	return false
}

func GetValidMessage(msg string) string {
	illegalChars := MakeRange(0, 32)
	illegalChars = append(illegalChars, 127)

	chars := []byte(msg)

	// check if string has at least one legal character:
	isValid := false
	for char := range chars {
		if !isInRange(illegalChars, char) {
			isValid = true
			break
		}
	}

	if !isValid {
		return ""
	}

	// remove null character
	var validChars []byte
	for char := range chars {
		switch char {
		case 0:
			validChars = append(validChars, 32)
		default:
			validChars = append(validChars, uint8(char))
		}

	}

	return string(validChars)
}

func IsSpaceString(txt string) bool {
	trimTxt := strings.TrimSpace(txt)
	return len(trimTxt) == 0
}

func RemoveNullChars(txt string) string {
	byteTxt := bytes.Trim([]byte(txt), "\x00")
	return string(byteTxt)
}

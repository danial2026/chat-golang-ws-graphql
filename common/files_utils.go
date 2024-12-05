package common

import "strings"

func IsValidMimeType(mimeString, majorType, minorType string) bool {
	mimeSlice := strings.Split(mimeString, "/")
	if len(mimeSlice) != 2 && mimeSlice[0] != majorType {
		return false
	}
	if minorType != "" {
		if mimeSlice[1] != minorType {
			return false
		}
	}
	return true
}

package common

import "github.com/google/uuid"

func GetUUID() string {
	id := uuid.New()
	return id.String()
}

func GetAndValidateUUID(input string) (string, error) {
	v, err := uuid.Parse(input)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

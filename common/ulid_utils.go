package common

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func GetULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

func GetAndValidateULID(input string) (string, error) {
	id, err := ulid.ParseStrict(input)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

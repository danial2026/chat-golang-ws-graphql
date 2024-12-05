package common

import "time"

func Now() int {
	return int(time.Now().Unix())
}

package common

import (
	"time"
)

func GetNowTimeStampS() int64 {
	t := time.Now()
	return int64(t.UTC().Second())
}

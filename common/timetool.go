package common

import (
	"time"
)

func GetNowTimeStampS() int64 {
	return time.Now().Unix()
}

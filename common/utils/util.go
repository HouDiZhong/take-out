package utils

import (
	"strconv"
	"take-out/common/enum"
	"time"
)

func Now() string {
	return time.Now().Format(enum.TimeLayout)
}

func TimeStamp() int64 {
	return time.Now().Unix()
}

func TimeStampStr() string {
	return strconv.FormatInt(TimeStamp(), 10)
}

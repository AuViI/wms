package forecast

import (
	"fmt"
	"time"
)

func nowString() string {
	now := time.Now().Local()
	return tString(now.Unix())
}

func tString(dt int64) string {
	now := time.Unix(dt, 0)
	return fmt.Sprintf("%02d.%02d.%04d %02d:%02d:%02d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
}

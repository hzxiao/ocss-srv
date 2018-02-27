package tools

import "time"

func NowMillisecond() int64 {
	return time.Now().Local().UnixNano() / 1e6
}

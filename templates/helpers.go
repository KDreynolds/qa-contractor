package templates

import "time"

func yearSpan() int {
	return time.Now().Year() - 2017
}

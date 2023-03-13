package geecache

import "log"

const debug = true

func DPrintf(format string, a ...interface{}) {
	if debug {
		log.Printf(format, a...)
	}
	return
}

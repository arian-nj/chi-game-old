package utils

import (
	"fmt"
	"runtime/debug"
)

func RunBackgroundTask(task func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("background task panic: %v\n%s", r, debug.Stack())
			}
		}()
		task()
	}()
}

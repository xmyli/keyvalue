package raft

import (
	"math/rand"
	"time"
)

func sleepFor(min int, max int) {
	duration := int64(min)
	if max-min > 0 {
		duration += rand.Int63() % int64(max-min)
	}
	time.Sleep(time.Duration(duration) * time.Millisecond)
}

func isMajorityGreaterOrEqual(x int, list []int) bool {
	count := 0
	for _, n := range list {
		if n >= x {
			count++
		}
	}
	return count+1 > len(list)/2
}

package main

import "time"

type timeSlice []time.Time

// Len is used by sort for sorting timestamps.
func (ts timeSlice) Len() int {
	return len(ts)
}

// Less is used by sort for sorting timestamps.
func (ts timeSlice) Less(i, j int) bool {
	return ts[i].Before(ts[j])
}

// Swap is used by sort for sorting timestamps.
func (ts timeSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

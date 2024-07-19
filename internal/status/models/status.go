package models

import "time"

type ApplicationUpTime struct {
	Requests      int64     `json:"requests"`
	StartedAt     time.Time `json:"startedAt"`
	LifetimeSecs  int64     `json:"lifetimeSecs"`
	Status        string    `json:"status"`
	NumGoroutine  int       `json:"numGoroutine"`
	NumGCCycles   uint32    `json:"numGCCycles"`
	AllocMiB      float64   `json:"allocMiB"`
	TotalAllocMiB float64   `json:"totalAllocMiB"`
}

package common

import (
	"runtime"
	"sync"
	"time"
)

type SingletonApplicationUpTime interface {
	AddRequest()
	RequestCount() int64
	StartedAt() time.Time
	MemStats() *runtime.MemStats
}

var singletonApplicationUpTimeInstance *applicationUpTime

func init() {
	singletonApplicationUpTimeInstance = &applicationUpTime{
		startedAt: time.Now().UTC(),
	}
}

func GetSingletonApplicationUpTime() SingletonApplicationUpTime {
	return singletonApplicationUpTimeInstance
}

type applicationUpTime struct {
	mutex        sync.Mutex
	requestCount int64
	startedAt    time.Time
}

func (appUpTime *applicationUpTime) AddRequest() {
	appUpTime.mutex.Lock()
	appUpTime.requestCount++
	appUpTime.mutex.Unlock()
}

func (appUpTime *applicationUpTime) RequestCount() int64 {
	appUpTime.mutex.Lock()
	defer appUpTime.mutex.Unlock()
	return appUpTime.requestCount
}

func (appUpTime *applicationUpTime) StartedAt() time.Time {
	return appUpTime.startedAt
}

func (appUpTime *applicationUpTime) MemStats() *runtime.MemStats {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	return memStats
}

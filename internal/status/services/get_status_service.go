package services

import (
	"context"
	"runtime"
	"time"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/common/persistence"
	"github.com/gustavlsouz/documents-service/internal/status/models"
)

type GetStatusService interface {
	Execute(ctx context.Context) *models.ApplicationUpTime
}

func NewGetStatusService() GetStatusService {
	return &getStatusService{}
}

type getStatusService struct{}

func (service *getStatusService) Execute(ctx context.Context) *models.ApplicationUpTime {

	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appUpTime := common.GetSingletonApplicationUpTime()
	memStats := appUpTime.MemStats()

	startedAt := appUpTime.StartedAt()
	result := &models.ApplicationUpTime{
		Requests:      appUpTime.RequestCount(),
		StartedAt:     startedAt,
		LifetimeSecs:  int64(time.Since(startedAt).Seconds()),
		Status:        "UP",
		NumGoroutine:  runtime.NumGoroutine(),
		NumGCCycles:   memStats.NumForcedGC,
		AllocMiB:      common.ToMiB(memStats.Alloc),
		TotalAllocMiB: common.ToMiB(memStats.TotalAlloc),
	}

	err := persistence.GetPersistenceInstance().Database().PingContext(timeoutContext)

	if err != nil {
		result.Status = "DOWN"
	}

	return result
}

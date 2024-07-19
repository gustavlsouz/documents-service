package operations

import (
	"time"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
)

type DocumentUpdater struct {
	payload *models.Document
}

func (inserter *DocumentUpdater) TableName() string {
	return "document"
}

func (inserter *DocumentUpdater) Fields() []interface{} {
	if inserter.payload.IsBlocked.Valid {
		return []interface{}{
			inserter.payload.Id,
			inserter.payload.Value,
			inserter.payload.UpdatedAt,
			inserter.payload.IsBlocked,
		}
	}

	return []interface{}{
		inserter.payload.Id,
		inserter.payload.Value,
		inserter.payload.UpdatedAt,
	}
}

func (inserter *DocumentUpdater) Statement() string {
	if inserter.payload.IsBlocked.Valid {
		return `update document set "value" = $2, "updatedAt" = $3, "isBlocked" = $4 where id = $1`
	}
	return `update document set "value" = $2, "updatedAt" = $3 where id = $1`
}

func (inserter *DocumentUpdater) Data() interface{} {
	return inserter.payload
}

func NewDocumentUpdaterCreator() *DocumentUpdaterCreator {
	return &DocumentUpdaterCreator{}
}

type DocumentUpdaterCreator struct{}

func (creator *DocumentUpdaterCreator) Create(payload *models.DocumentPayload) common.WriteOperation {
	return &DocumentUpdater{
		payload: &models.Document{
			Id:        payload.Id,
			Value:     payload.Value,
			UpdatedAt: time.Now().UTC(),
			IsBlocked: payload.IsBlocked,
		},
	}
}

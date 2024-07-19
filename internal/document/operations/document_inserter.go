package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
)

type DocumentInserter struct {
	payload *models.Document
}

func (inserter *DocumentInserter) TableName() string {
	return "document"
}

func (inserter *DocumentInserter) Fields() []interface{} {
	return []interface{}{
		inserter.payload.Id,
		inserter.payload.Type,
		inserter.payload.Value,
		inserter.payload.CreatedAt,
		inserter.payload.UpdatedAt,
		inserter.payload.IsBlocked,
	}
}

func (inserter *DocumentInserter) Statement() string {
	return `insert into "document" 
	("id", "type", "value", "createdAt", "updatedAt", "isBlocked") 
	values ($1, $2, $3, $4, $5, $6)`
}

func (inserter *DocumentInserter) Data() interface{} {
	return inserter.payload
}

func NewDocumentInserterCreator() *DocumentInserterCreator {
	return &DocumentInserterCreator{}
}

type DocumentInserterCreator struct{}

func (creator *DocumentInserterCreator) Create(payload *models.DocumentPayload) common.WriteOperation {
	return &DocumentInserter{
		payload: &models.Document{
			Id:        uuid.NewString(),
			Type:      payload.Type,
			Value:     payload.Value,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			IsBlocked: payload.IsBlocked,
		},
	}
}

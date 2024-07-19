package operations

import (
	"fmt"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
)

type DocumentReader struct {
	payload *models.Document
}

func (reader *DocumentReader) TableName() string {
	return "document"
}

func (reader *DocumentReader) Args() []interface{} {
	if reader.payload.Id != "" {
		return []interface{}{reader.payload.Id}
	}

	if reader.payload.Value != "" {
		return []interface{}{reader.payload.Value}
	}

	return []interface{}{}
}

func (reader *DocumentReader) Query(pagination common.Pagination) string {
	skip := (pagination.Page() - 1) * pagination.Size()
	if reader.payload.Id != "" {
		return `select id,
		"type",
		"value",
		"createdAt", 
		"updatedAt", 
		"isBlocked" 
		from document where id = $1`
	}

	if reader.payload.Value != "" {
		return fmt.Sprintf(`select id,
		"type",
		"value",
		"createdAt", 
		"updatedAt", 
		"isBlocked"
		from document where value like  '%%' || $1 || '%%' order by "updatedAt" desc limit %d offset %d`,
			pagination.Size(),
			skip)
	}

	// no pagination to simplify
	return fmt.Sprintf(`select id,
	"type",
	"value",
	"createdAt", 
	"updatedAt", 
	"isBlocked"
	from document order by "updatedAt" desc limit %d offset %d`, pagination.Size(), skip)
}

func NewDocumentReaderCreator() *DocumentReaderCreator {
	return &DocumentReaderCreator{}
}

type DocumentReaderCreator struct{}

func (creator *DocumentReaderCreator) Create(payload *models.Document) common.ReadOperation {
	return &DocumentReader{
		payload: payload,
	}
}

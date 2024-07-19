package operations

import (
	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
)

type DocumentRemover struct {
	payload *models.Document
}

func (remover *DocumentRemover) TableName() string {
	return "document"
}

func (remover *DocumentRemover) Fields() []interface{} {
	return []interface{}{remover.payload.Id}
}

func (remover *DocumentRemover) Statement() string {
	return "delete from document where id = $1"
}

func (remover *DocumentRemover) Data() interface{} {
	return remover.payload
}

func NewDocumentRemoverCreator() *documentRemoverCreator {
	return &documentRemoverCreator{}
}

type documentRemoverCreator struct {
}

func (creator *documentRemoverCreator) Create(payload *models.Document) common.WriteOperation {
	return &DocumentRemover{
		payload: payload,
	}
}

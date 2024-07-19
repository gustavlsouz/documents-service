package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
	"github.com/gustavlsouz/documents-service/internal/document/operations"
	"github.com/gustavlsouz/documents-service/internal/wrappers"
)

type ValidateDocumentService interface {
	Execute(ctx context.Context, document *models.DocumentPayload) (any, error)
}

func NewValidateDocumentService() ValidateDocumentService {
	return &validateDocumentService{
		documentValidator: wrappers.NewDocumentValidator(),
		readerRepository: common.NewReaderRepository[models.Document, models.Document](
			operations.NewDocumentReaderCreator(),
		),
		documentFormatter: wrappers.NewDocumentFormatter(),
	}
}

type validateDocumentService struct {
	documentValidator models.DocumentValidator
	readerRepository  common.ReaderRepository[models.Document, models.Document]
	documentFormatter models.DocumentFormatter
}

var ErrorInvalidDocument = errors.New("invalid document")

func (service *validateDocumentService) Execute(ctx context.Context, document *models.DocumentPayload) (any, error) {
	log.Println(document.Type, document.Value)

	documents, err := service.readerRepository.Read(ctx, &models.Document{Id: document.Id}, common.NewPagination())

	if err != nil {
		return nil, fmt.Errorf("error to read document by id: %w", err)
	}

	if len(documents) == 0 {
		return nil, errors.New("document does not exists")
	}

	currentDocument := documents[0]

	document.Type = currentDocument.Type
	document.Value = service.documentFormatter.CleanPad(currentDocument.Type, document.Value)

	err = document.Validate(service.documentValidator)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("is invalid document: %w", err)
	}

	log.Println("success to validate document")
	return nil, nil
}

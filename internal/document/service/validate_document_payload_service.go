package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gustavlsouz/documents-service/internal/document/models"
	"github.com/gustavlsouz/documents-service/internal/wrappers"
)

type ValidateDocumentPayloadService interface {
	Execute(ctx context.Context, documentPayload *models.DocumentPayload) (any, error)
}

func NewValidateDocumentPayloadService() ValidateDocumentPayloadService {
	return &validateDocumentPayloadService{
		documentValidator: wrappers.NewDocumentValidator(),
		documentFormatter: wrappers.NewDocumentFormatter(),
	}
}

type validateDocumentPayloadService struct {
	documentValidator models.DocumentValidator
	documentFormatter models.DocumentFormatter
}

func (service *validateDocumentPayloadService) Execute(ctx context.Context, documentPayload *models.DocumentPayload) (any, error) {

	log.Println(documentPayload.Type, documentPayload.Value)

	documentPayload.Value = service.documentFormatter.
		CleanPad(documentPayload.Type, documentPayload.Value)

	err := documentPayload.Validate(service.documentValidator)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("is invalid document: %w", err)
	}

	log.Println("success to validate document")
	return nil, nil
}

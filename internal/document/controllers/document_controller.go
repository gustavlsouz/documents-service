package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/document/models"
	"github.com/gustavlsouz/documents-service/internal/document/service"
	"github.com/gustavlsouz/documents-service/internal/wrappers"
)

type documentQueryCreator struct{}

func NewDocumentQueryCreator() *documentQueryCreator {
	return &documentQueryCreator{}
}

func (creator *documentQueryCreator) Create(request *http.Request) (*models.Document, error) {

	document := &models.Document{
		Id:    request.URL.Query().Get("id"),
		Value: request.URL.Query().Get("value"),
	}

	return document, nil
}

type documentDeleteCriteriaCreator struct{}

func NewDocumentDeleteCriteriaCreator() *documentDeleteCriteriaCreator {
	return &documentDeleteCriteriaCreator{}
}

func (creator *documentDeleteCriteriaCreator) Create(request *http.Request) (*models.Document, error) {
	documentId := request.URL.Query().Get("id")
	return &models.Document{Id: documentId}, nil
}

type DocumentResponse struct {
	Id        string              `json:"id,omitempty"`
	Type      models.DocumentType `json:"type,omitempty"`
	Value     string              `json:"value,omitempty"`
	CreatedAt string              `json:"createdAt,omitempty"`
	UpdatedAt string              `json:"updatedAt,omitempty"`
	IsBlocked common.JsonNullBool `json:"isBlocked"`
}

type documentReadFormatter struct {
	documentFormatter models.DocumentFormatter
}

func NewDocumentReadFormatter() common.ReadFormatter {
	return &documentReadFormatter{
		documentFormatter: wrappers.NewDocumentFormatter(),
	}
}

func (formatter *documentReadFormatter) adjustTime(location *time.Location, date time.Time) time.Time {
	if location == nil {
		return date
	}
	return date.In(location)
}

func (formatter *documentReadFormatter) formatTime(date time.Time) string {
	return fmt.Sprintf("%s %s", date.Format("02/01/2006"), date.Format("15:04:05"))
}

func (formatter *documentReadFormatter) FormatAll(request *http.Request, data interface{}) (interface{}, error) {
	list, ok := data.([]models.Document)
	if !ok {
		return data, nil
	}

	timezone := request.Header.Get("X-Timezone")

	log.Println("timezone", timezone)

	location, err := time.LoadLocation(timezone)

	if err != nil {
		return nil, err
	}

	formatteds := make([]DocumentResponse, 0)

	for _, item := range list {
		response := DocumentResponse{
			Id:        item.Id,
			Type:      item.Type,
			Value:     formatter.documentFormatter.Format(item.Type, item.Value),
			CreatedAt: formatter.formatTime(formatter.adjustTime(location, item.CreatedAt)),
			UpdatedAt: formatter.formatTime(formatter.adjustTime(location, item.UpdatedAt)),
			IsBlocked: item.IsBlocked,
		}
		formatteds = append(formatteds, response)
	}

	return formatteds, nil
}

func NewDocumentController(
	reader common.ReadOperationCreator[models.Document],
	inserterCreator common.WriteOperationCreator[models.DocumentPayload],
	deleterCreator common.WriteOperationCreator[models.Document],
	updaterCreator common.WriteOperationCreator[models.DocumentPayload],
) common.CrudController[models.DocumentPayload, models.Document, models.Document] {
	return common.NewCrudController(
		common.NewReaderWithOperationCreator[models.Document, models.Document](reader),
		common.NewWriterWithOperationCreatorWithCustomService(service.NewValidateDocumentPayloadService(), inserterCreator),
		common.NewWriterWithOperationCreatorWithCustomService(service.NewValidateDocumentService(), updaterCreator),
		common.NewWriterWithOperationCreator(deleterCreator),
		NewDocumentQueryCreator(),
		NewDocumentDeleteCriteriaCreator(),
		NewDocumentReadFormatter(),
	)
}

package integration

import (
	"errors"
	"log"
	"testing"

	"github.com/gustavlsouz/documents-service/internal/document/controllers"
	"github.com/gustavlsouz/documents-service/internal/document/models"
	"github.com/gustavlsouz/documents-service/internal/wrappers"
)

var validator models.DocumentValidator
var formatter models.DocumentFormatter

func init() {
	validator = wrappers.NewDocumentValidator()
	formatter = wrappers.NewDocumentFormatter()
}

func NewRandomValue() string {
	cpfText := GenerateCpf()
	if !validator.IsCPF(cpfText) {
		log.Panic(errors.New("it's not a valid document"))
	}

	return cpfText
}

type DocumentUpdater struct{}

func (documentUpdater *DocumentUpdater) Update(document models.Document) models.Document {
	document.Value = NewRandomValue()
	return document
}

type DocumentComparer struct{}

func (documentComparer *DocumentComparer) Compare(documentX models.Document, response controllers.DocumentResponse) bool {
	return documentX.Value == formatter.Clean(response.Value)
}

type DocumentIdentifier struct{}

func (documentIdentifier *DocumentIdentifier) Identify(document models.Document) string {
	if document.Id == "" {
		return document.Value
	}
	return document.Id
}

func NewMockedDocument() models.Document {
	return models.Document{
		Type:  "CPF",
		Value: NewRandomValue(),
	}
}

func TestDocument(t *testing.T) {

	httpUrl := "http://localhost:8080/api/document"
	documentTest := NewMockedDocument()

	crudTesting := NewCrudTesting[models.Document](
		"document",
		&DocumentUpdater{},
		&DocumentComparer{},
		&DocumentIdentifier{},
	)

	crudTesting.Execute(t, httpUrl, documentTest)

}

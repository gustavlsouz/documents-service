package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gustavlsouz/documents-service/internal/common"
)

type DocumentType string

const (
	CPF  DocumentType = "CPF"
	CNPJ DocumentType = "CNPJ"
)

var documentTypes = []DocumentType{CPF, CNPJ}

type DocumentPayload struct {
	Id        string              `json:"id,omitempty"`
	Type      DocumentType        `json:"type,omitempty"`
	Value     string              `json:"value"`
	IsBlocked common.JsonNullBool `json:"isBlocked"`
}

var ErrorInvalidCNPJ = errors.New("invalid CNPJ")
var ErrorInvalidCPF = errors.New("invalid CPF")
var ErrorInvalidDocumentType = errors.New("invalid document type")

func (payload *DocumentPayload) GetType() DocumentType {
	return payload.Type
}

func (payload *DocumentPayload) GetValue() string {
	return payload.Value
}

func (model *DocumentPayload) Validate(validator DocumentValidator) error {
	return Validate(validator, model)
}

type Document struct {
	Id        string              `json:"id,omitempty"`
	Type      DocumentType        `json:"type,omitempty"`
	Value     string              `json:"value,omitempty"`
	CreatedAt time.Time           `json:"createdAt,omitempty"`
	UpdatedAt time.Time           `json:"updatedAt,omitempty"`
	IsBlocked common.JsonNullBool `json:"isBlocked"`
}

func (payload *Document) GetType() DocumentType {
	return payload.Type
}

func (payload *Document) GetValue() string {
	return payload.Value
}

func (model *Document) Validate(validator DocumentValidator) error {
	return Validate(validator, model)
}

type ValidatableDocument interface {
	GetType() DocumentType
	GetValue() string
}

func Validate(validator DocumentValidator, validatable ValidatableDocument) error {
	value := validatable.GetValue()
	if value == "" {
		return errors.New("empty document")
	}
	isValidType := IsValidDocumentType(validatable.GetType())
	if !isValidType {
		return fmt.Errorf("'%s' is not valid: %w", validatable.GetType(), ErrorInvalidDocumentType)
	}

	switch validatable.GetType() {
	case CNPJ:
		return common.RaiseWhenNok(validator.IsCNPJ(FormatByType(validatable.GetType(), validatable.GetValue())), ErrorInvalidCNPJ)
	case CPF:
		return common.RaiseWhenNok(validator.IsCPF(FormatByType(validatable.GetType(), validatable.GetValue())), ErrorInvalidCPF)
	default:
		return fmt.Errorf("'%s' is not valid: %w", validatable.GetType(), ErrorInvalidDocumentType)
	}
}

func IsValidDocumentType(modelType DocumentType) bool {
	for _, typeItem := range documentTypes {
		if modelType == typeItem {
			return true
		}
	}
	return false
}

func FormatByType(modelType DocumentType, value string) string {
	log.Println("type", modelType)
	switch modelType {
	case CNPJ:
		cnpj := fmt.Sprintf("%014s", value)
		return cnpj
	case CPF:
		cpf := fmt.Sprintf("%011s", value)
		return cpf
	default:
		return ""
	}
}

type DocumentValidator interface {
	IsCNPJ(cnpj string) bool
	IsCPF(cpf string) bool
}
type DocumentFormatter interface {
	Format(documentType DocumentType, documentValue string) string
	Clean(string) string
	CleanPad(documentType DocumentType, document string) string
}

package wrappers

import (
	"fmt"

	"github.com/gustavlsouz/documents-service/internal/document/models"
	"github.com/klassmann/cpfcnpj"
)

func NewDocumentValidator() models.DocumentValidator {
	return &documentValidatorWrapper{}
}

type documentValidatorWrapper struct{}

func (validator *documentValidatorWrapper) IsCNPJ(cnpj string) bool {
	return cpfcnpj.ValidateCNPJ(cnpj)
}

func (validator *documentValidatorWrapper) IsCPF(cpf string) bool {
	return cpfcnpj.ValidateCPF(cpf)
}

func NewDocumentFormatter() models.DocumentFormatter {
	return &documentFormatterWrapper{}
}

type documentFormatterWrapper struct{}

func (formatterUtil *documentFormatterWrapper) Format(documentType models.DocumentType, documentValue string) string {
	switch documentType {
	case models.CNPJ:
		return formatterUtil.FormatCNPJ(documentValue)
	case models.CPF:
		return formatterUtil.FormatCPF(documentValue)
	default:
		return ""
	}
}

func (formatterUtil *documentFormatterWrapper) FormatCNPJ(cnpj string) string {
	c := cpfcnpj.NewCNPJ(fmt.Sprintf("%014s", cnpj))
	return c.String()
}

func (formatterUtil *documentFormatterWrapper) FormatCPF(cpf string) string {
	c := cpfcnpj.NewCPF(fmt.Sprintf("%011s", cpf))
	return c.String()
}

func (formatterUtil *documentFormatterWrapper) CleanPad(documentType models.DocumentType, document string) string {

	switch documentType {
	case models.CNPJ:
		return cpfcnpj.Clean(formatterUtil.FormatCNPJ(document))
	case models.CPF:
		return cpfcnpj.Clean(formatterUtil.FormatCPF(document))
	default:
		return ""
	}
}

func (formatterUtil *documentFormatterWrapper) Clean(document string) string {
	return cpfcnpj.Clean(document)
}

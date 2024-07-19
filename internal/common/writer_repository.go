package common

import (
	"context"
	"fmt"
	"log"

	"github.com/gustavlsouz/documents-service/internal/common/persistence"
)

type WriterRepository[T any] interface {
	Write(context.Context, *T) (any, error)
}

func NewWriterRepository[T any](operationCreator WriteOperationCreator[T]) WriterRepository[T] {
	return &writerRepository[T]{
		operationCreator: operationCreator,
		persistence:      persistence.GetPersistenceInstance(),
	}
}

type writerRepository[T any] struct {
	operationCreator WriteOperationCreator[T]
	persistence      persistence.Persistence
}

func (wRepository *writerRepository[T]) Write(ctx context.Context, model *T) (any, error) {
	operation := wRepository.operationCreator.Create(model)
	statment := operation.Statement()
	fields := operation.Fields()
	log.Println("changing:", operation.TableName(), "statment:", statment, "fields:", fields)
	_, err := wRepository.persistence.Database().ExecContext(ctx, statment, fields...)

	if err != nil {
		return nil, fmt.Errorf("error to execute statment: %w", err)
	}
	return operation.Data(), nil
}

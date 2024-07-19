package common

import (
	"context"
)

type ReaderService[Q any, T any] interface {
	Execute(ctx context.Context, model *Q, pagination Pagination) ([]T, error)
}

func NewReaderService[Q any, T any](repository ReaderRepository[Q, T]) ReaderService[Q, T] {
	return &readerService[Q, T]{
		writerRepository: repository,
	}
}

type readerService[Q any, T any] struct {
	writerRepository ReaderRepository[Q, T]
}

func (operator *readerService[Q, T]) Execute(ctx context.Context, model *Q, pagination Pagination) ([]T, error) {
	return operator.writerRepository.Read(ctx, model, pagination)
}

func NewReaderWithOperationCreator[Q any, T any](operationCreator ReadOperationCreator[Q]) ReaderService[Q, T] {
	return NewReaderService(
		NewReaderRepository[Q, T](operationCreator),
	)
}

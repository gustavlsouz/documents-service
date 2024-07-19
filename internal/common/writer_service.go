package common

import (
	"context"
)

type WriterService[T any] interface {
	Execute(ctx context.Context, model *T) (any, error)
}

func NewWriterService[T any](repository WriterRepository[T]) WriterService[T] {
	return &writerService[T]{
		writerRepository: repository,
	}
}

func NewWriterServiceWithCustomService[T any](customService WriterService[T], repository WriterRepository[T]) WriterService[T] {
	return &writerService[T]{
		customService:    customService,
		writerRepository: repository,
	}
}

type writerService[T any] struct {
	writerRepository WriterRepository[T]
	customService    WriterService[T]
}

func (operator *writerService[T]) Execute(ctx context.Context, model *T) (any, error) {
	if operator.customService == nil {
		return operator.writerRepository.Write(ctx, model)
	}

	customResult, err := operator.customService.Execute(ctx, model)
	if err != nil {
		return nil, err
	}
	if customResult != nil {
		return customResult, nil
	}
	return operator.writerRepository.Write(ctx, model)
}

func NewWriterWithOperationCreator[T any](operationCreator WriteOperationCreator[T]) WriterService[T] {
	return NewWriterService(
		NewWriterRepository(operationCreator),
	)
}

func NewWriterWithOperationCreatorWithCustomService[T any](
	customService WriterService[T],
	operationCreator WriteOperationCreator[T]) WriterService[T] {
	return NewWriterServiceWithCustomService(
		customService,
		NewWriterRepository(operationCreator),
	)
}

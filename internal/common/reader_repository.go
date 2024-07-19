package common

import (
	"context"
	"fmt"
	"log"

	"github.com/gustavlsouz/documents-service/internal/common/persistence"
)

type ReaderRepository[Q any, T any] interface {
	Read(context.Context, *Q, Pagination) ([]T, error)
}

func NewReaderRepository[Q any, T any](operationCreator ReadOperationCreator[Q]) ReaderRepository[Q, T] {
	return &readerRepository[Q, T]{
		operationCreator: operationCreator,
		persistence:      persistence.GetPersistenceInstance(),
	}
}

type readerRepository[Q any, T any] struct {
	operationCreator ReadOperationCreator[Q]
	persistence      persistence.Persistence
}

func (rRepository *readerRepository[Q, T]) Read(ctx context.Context, model *Q, pagination Pagination) ([]T, error) {
	operation := rRepository.operationCreator.Create(model)
	query := operation.Query(pagination)
	args := operation.Args()
	log.Println("reading:", operation.TableName(), "query:", query, "args:", args)
	rows, err := rRepository.persistence.Database().QueryContext(ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("error to query: %w", err)
	}

	list := make([]T, 0)

	for rows.Next() {
		var item T
		err = persistence.GetPersistenceInstance().ScanStruct(rows, &item)

		if err != nil {
			return nil, fmt.Errorf("error to scan struct: %w", err)
		}

		list = append(list, item)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("error to close rows: %w", err)
	}

	return list, nil
}

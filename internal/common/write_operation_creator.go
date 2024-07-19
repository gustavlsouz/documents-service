package common

type WriteOperationCreator[T any] interface {
	Create(*T) WriteOperation
}

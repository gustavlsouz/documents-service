package common

type ReadOperationCreator[Q any] interface {
	Create(*Q) ReadOperation
}

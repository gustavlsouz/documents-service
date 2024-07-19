package common

type ReadOperation interface {
	TableName() string
	Args() []interface{}
	Query(Pagination) string
}

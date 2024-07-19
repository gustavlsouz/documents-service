package common

type WriteOperation interface {
	TableName() string
	Fields() []interface{}
	Statement() string
	Data() interface{}
}

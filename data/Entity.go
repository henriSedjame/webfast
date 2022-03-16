package data

//Entity represents a sql entity
type Entity interface {
	TableName() string
}

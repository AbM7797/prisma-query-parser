package parser

// TypeMapper defines the interface for managing project-specific types
type TypeMapper interface {
	GetDBInstance(key string) interface{}
	GetWhereClause(key string) interface{}
	GetOrderByClause(key string) interface{}
	GetDomainType(key string) interface{}
	GetASC() interface{}
	GetDESC() interface{}
	GetMode(key string) interface{}
	GetValue(key string) interface{}
}

// Filter represents the basic structure of filters
type Filter map[string]interface{}

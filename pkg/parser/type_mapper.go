package parser

// TypeMapper définit l'interface pour la gestion des types spécifiques au projet
type TypeMapper interface {
	GetDBInstance(key string) interface{}
	GetWhereClause(key string) interface{}
	GetOrderByClause(key string) interface{}
	GetDomainType(key string) interface{}
	GetASC() interface{}
	GetDESC() interface{}
}

// Filter représente la structure de base des filtres
type Filter map[string]interface{}

# Prisma Query Parser

Prisma Query Parser is a Go library that helps parse and process query parameters into a structured filter map, making it easier to work with Prisma databases.

## Features

* Process query parameters into a structured filter map
* Handle WHERE, ORDER BY, SKIP, and TAKE clauses
* Dynamic mapping of filter operators to Prisma methods
* Support for logical operations (AND, OR, NOT)
* Type mapping and conversion for project-specific data types

## Getting Started

1. Install the library using Go modules: `go get github.com/AbM7797/prisma-query-parser/`
2. Import the library in your Go project: `import "github.com/AbM7797/prisma-query-parser/"`
3. Use the `ProcessArgs` function to parse query parameters into a filter map
4. Use the `BuildOrderFilters` and `BuildWhereFilters` functions to build and process filters for your Prisma database

## API Documentation

### ProcessArgs

* `func ProcessArgs(args url.Values) Filter`
	+ Parse query parameters into a structured filter map
	+ Returns a `Filter` map containing the parsed query parameters

#### Example of code
```
func (h *PostHandler) Find(c *gin.Context) {
	filters := parser.ProcessArgs(c.Request.URL.Query())
	posts, err := h.useCase.Find(c.Request.Context(), filters)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, posts)
}
```

### BuildOrderFilters

* `func BuildOrderFilters[U any](filters Filter, customFilters []U, tm TypeMapper, mapper string) []U`
	+ Build order filters for a Prisma database
	+ Returns a slice of `U` containing the built order filters

### WhereMapper

* `func WhereMapper[T any, K any, C any](model K, filter Filter, returnType T, where C, tm TypeMapper) []T`
	+ Map a filter to a Prisma method
	+ Returns a slice of `T` containing the mapped filter

#### Example of code
```
func (r *postRepository) Find(ctx context.Context, filters parser.Filter, customWhereFilters []db.PostWhereParam, customOrderByFilter []db.PostOrderByParam) ([]db.PostModel, error) {
	whereFilters := parser.BuildWhereFilters(
		filters,
		customWhereFilters,
		r.tm,
		"post",
	)
	orderByFilters := parser.BuildOrderFilters(filters, customOrderByFilter, r.tm, "post")
	query := r.client.Post.FindMany(whereFilters...).
		OrderBy(orderByFilters...).
		With(db.Post.User.Fetch())

	if filters["take"] != nil {
		query = query.Take(filters["take"].(int))
	}

	if filters["skip"] != nil {
		query = query.Skip(filters["skip"].(int))
	}

	posts, err := query.Exec(ctx)

	if err != nil {
		return nil, err
	}

	return posts, nil
}
```

## Type Mapping

The library uses a `TypeMapper` interface to manage project-specific data types. You can implement this interface to provide custom type mapping for your project.

#### Example of code
```
package types

import (
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
)

type TypeMapper struct {
	types map[string]interface{}
}

func NewTypeMapper() *TypeMapper {
	var userWhereParam db.UserWhereParam
	userOrderByParam := new(db.UserOrderByParam)

	var postWhereParam db.PostWhereParam
	postOrderByParam := new(db.PostOrderByParam)

	var commentWhereParam db.CommentWhereParam
	commentOrderByParam := new(db.CommentOrderByParam)

	return &TypeMapper{
		types: map[string]interface{}{
			//ASC
			"asc":  db.ASC,
			"desc": db.DESC,

			// user
			"user_where":    userWhereParam,
			"user_order_by": userOrderByParam,
			"user_db":       db.User,
			"user_domain":   &domain.User{},

			// post
			"post_where":    postWhereParam,
			"post_order_by": postOrderByParam,
			"post_db":       db.Post,
			"post_domain":   &domain.Post{},

			// comment
			"comment_where":    commentWhereParam,
			"comment_order_by": commentOrderByParam,
			"comment_db":       db.Comment,
			"comment_domain":   &domain.Comment{},
		},
	}
}

func (t *TypeMapper) GetDBInstance(key string) interface{} {
	return t.types[key+"_db"]
}

func (t *TypeMapper) GetWhereClause(key string) interface{} {
	return t.types[key+"_where"]
}

func (t *TypeMapper) GetDomainType(key string) interface{} {
	return t.types[key+"_domain"]
}

func (t *TypeMapper) GetOrderByClause(key string) interface{} {
	return t.types[key+"_order_by"]
}

func (t *TypeMapper) GetASC() interface{} {
	return t.types["asc"]
}

func (t *TypeMapper) GetDESC() interface{} {
	return t.types["desc"]
}
```

## Contributing

We welcome contributions to the Prisma Query Parser library. If you'd like to contribute, please fork the repository and submit a pull request with your changes.

## License

Prisma Query Parser is licensed under the MIT License.

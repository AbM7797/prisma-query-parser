package types

import (
	"github.com/AbM7797/prisma-query-parser/examples/app/internal/domain"
	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
)

type TypeMapper struct {
	types map[string]interface{}
}

func NewTypeMapper() *TypeMapper {
	userWhereParam := new(db.UserWhereParam)
	userOrderByParam := new(db.UserOrderByParam)

	var postWhereParam db.PostWhereParam
	//postWhereParam := new(db.PostWhereParam)
	postOrderByParam := new(db.PostOrderByParam)

	commentWhereParam := new(db.CommentWhereParam)
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

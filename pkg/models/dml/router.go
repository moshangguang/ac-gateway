package dml

import (
	"ac-gateway/pkg/models/ddl"
	"context"
)

type RouterModel interface {
	GetAll(ctx context.Context) (list []ddl.Router, err error)
	GetById(ctx context.Context, id int64) (router ddl.Router, exists bool, err error)
	Delete(ctx context.Context, id int64) error
	Save(ctx context.Context, route ddl.Router) (ddl.Router, error)
}

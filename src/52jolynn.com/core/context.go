package core

import (
		"database/sql"
)

type Context interface {
	Datasource() *sql.DB
}

type context struct {
	datasource *sql.DB
}

func (ctx *context) Datasource() *sql.DB {
	return ctx.datasource
}

func NewContext(db *sql.DB) Context {
	return &context{datasource: db}
}

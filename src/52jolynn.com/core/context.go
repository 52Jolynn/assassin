package core

import (
	"github.com/kataras/golog"
	"database/sql"
)

type Context interface {
	RootLogger() *golog.Logger
	Datasource() *sql.DB
}

type context struct {
	rootLogger *golog.Logger
	datasource *sql.DB
}

func (ctx *context) RootLogger() *golog.Logger {
	return ctx.rootLogger
}

func (ctx *context) Datasource() *sql.DB {
	return ctx.datasource
}

func NewContext(rootLogger *golog.Logger, db *sql.DB) Context {
	return &context{rootLogger: rootLogger, datasource: db}
}

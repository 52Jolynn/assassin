package db

import (
	"fmt"
	"strings"
)

type SqlObject interface {
	ToSql() string
}

const (
	ColumnTypeAsterisk = "AsteriskColumn"
	ColumnTypeName     = "NameColumn"
	ColumnTypeValue    = "ValueColumn"
	ColumnTypeExpr     = "ExprColumn"
)

type IColumn interface {
	ToSql() string
	GetColumnType() string
}

type Column struct {
	Alias *string
	Index int
}

type AsteriskColumn struct {
}

func (c *AsteriskColumn) ToSql() string {
	return "*"
}

func (c *AsteriskColumn) GetColumnType() string {
	return ColumnTypeAsterisk
}

type NameColumn struct {
	Name string
	Column
}

func (c *NameColumn) ToSql() string {
	return fmt.Sprintf("%s as %s", c.Name, c.getAlias())
}

func (c *NameColumn) getAlias() string {
	if c.Alias == nil {
		return fmt.Sprintf("%s_%d", c.Name, c.Index)
	} else {
		return *c.Alias
	}
}

func (c *NameColumn) GetColumnType() string {
	return ColumnTypeName
}

type ValueColumn struct {
	Value interface{}
	Column
}

func (c *ValueColumn) ToSql() string {
	switch c.Value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d as %s", c.Value, c.getAlias())
	case string:
		return fmt.Sprintf("%s as %s", c.Value, c.getAlias())
	default:
		return fmt.Sprintf("%s as %s", c.Value, c.getAlias())
	}
}

func (c *ValueColumn) getAlias() string {
	if c.Alias == nil {
		return fmt.Sprintf("%s_%d", c.Value, c.Index)
	} else {
		return *c.Alias
	}
}

func (c *ValueColumn) GetColumnType() string {
	return ColumnTypeValue
}

type ExprColumn struct {
	expr Expression
	Column
}

func (c *ExprColumn) ToSql() string {
	return fmt.Sprintf("%s as %s", c.expr.ToSql(), c.getAlias())
}

func (c *ExprColumn) getAlias() string {
	if c.Alias == nil {
		return fmt.Sprintf("expr_%d", c.Index)
	} else {
		return *c.Alias
	}
}
func (c *ExprColumn) GetColumnType() string {
	return ColumnTypeExpr
}

const (
	ExprTypeValueExpr      = "ValueExpr"
	ExprTypeColumnExpr     = "ColumnExpr"
	ExprTypeUnaryExpr      = "UnaryExpr"
	ExprTypeBinaryExpr     = "BinaryExpr"
	ExprTypeNotExpr        = "NotExpr"
	ExprTypeAndExpr        = "AndExpr"
	ExprTypeOrExpr         = "OrExpr"
	ExprTypeBetweenExpr    = "BetweenExpr"
	ExprTypeNotBetweenExpr = "NotBetweenExpr"
)

type Expression interface {
	ToSql() string
	GetExprType() string
}

type ValueExpr struct {
	Value interface{}
}

func (e *ValueExpr) ToSql() string {
	switch e.Value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", e.Value)
	case string:
		return fmt.Sprintf("%s", e.Value)
	default:
		return fmt.Sprintf("%s", e.Value)
	}
}

func (e *ValueExpr) GetExprType() string {
	return ExprTypeValueExpr
}

type ColumnExpr struct {
	Alias string
	Name  string
}

func (e *ColumnExpr) ToSql() string {
	return fmt.Sprintf("%s.%s", e.Alias, e.Name)
}

func (e *ColumnExpr) GetExprType() string {
	return ExprTypeColumnExpr
}

type UnaryExpr struct {
	operator string
	expr     Expression
}

func (e *UnaryExpr) ToSql() string {
	return fmt.Sprintf("%s%s", e.operator, e.expr.ToSql())
}

func (e *UnaryExpr) GetExprType() string {
	return ExprTypeUnaryExpr
}

type NotExpr struct {
	expr Expression
}

func (e *NotExpr) ToSql() string {
	return "NOT " + e.expr.ToSql()
}

func (e *NotExpr) GetExprType() string {
	return ExprTypeNotExpr
}

//二元运算，如四则运算
type BinaryExpr struct {
	operator string
	left     Expression
	right    Expression
}

func (e *BinaryExpr) ToSql() string {
	return fmt.Sprintf("%s %s %s", e.left.ToSql(), e.operator, e.right.ToSql())
}

func (e *BinaryExpr) GetExprType() string {
	return ExprTypeBinaryExpr
}

type AndExpr struct {
	left  Expression
	right Expression
}

func (e *AndExpr) ToSql() string {
	return fmt.Sprintf("%s AND %s", e.left.ToSql(), e.right.ToSql())
}

func (e *AndExpr) GetExprType() string {
	return ExprTypeAndExpr
}

type OrExpr struct {
	AndExpr
}

func (e *OrExpr) ToSql() string {
	return fmt.Sprintf("%s OR %s", e.left.ToSql(), e.right.ToSql())
}

func (e *OrExpr) GetExprType() string {
	return ExprTypeOrExpr
}

type BetweenExpr struct {
	expr  Expression
	left  Expression
	right Expression
}

func (e *BetweenExpr) ToSql() string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", e.expr.ToSql(), e.left.ToSql(), e.right.ToSql())
}

func (e *BetweenExpr) GetExprType() string {
	return ExprTypeBetweenExpr
}

type NotBetweenExpr struct {
	BetweenExpr
}

func (e *NotBetweenExpr) ToSql() string {
	return fmt.Sprintf("%s NOT BETWEEN %s AND %s", e.expr.ToSql(), e.left.ToSql(), e.right.ToSql())
}

func (e *NotBetweenExpr) GetExprType() string {
	return ExprTypeNotBetweenExpr
}

type GroupByClause struct {
	columns []string
}

func (g *GroupByClause) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("GROUP BY ")
	for index, column := range g.columns {
		sql.WriteString(fmt.Sprintf(" %s", column))
		if index != len(g.columns)-1 {
			sql.WriteString(",")
		}
	}
	return sql.String()
}

type HavingClause struct {
	expr Expression
}

func (h *HavingClause) ToSql() string {
	return fmt.Sprintf("HAVING %s", h.expr.ToSql())
}

type OrderByClause struct {
	columns []OrderByColumn
}

type OrderByColumn struct {
	Name string
	Asc  bool
}

func (o *OrderByClause) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("ORDER BY ")
	for index, column := range o.columns {
		if column.Asc {
			sql.WriteString(fmt.Sprintf("%s ASC", column.Name))
		} else {
			sql.WriteString(fmt.Sprintf("%s DESC", column.Name))
		}
		if index != len(o.columns)-1 {
			sql.WriteString(",")
		}
	}
	return sql.String()
}

const (
	TableTypeTableName      = "TableName"
	TableTypeJoinTable      = "JoinTable"
	TableTypeCrossJoinTable = "CrossJoinTable"
	TableTypeSubquery       = "Subquery"
)

type Table interface {
	ToSql() string
	GetTableType() string
}

type TableName struct {
	Name  string
	Alias *string
}

func (t *TableName) ToSql() string {
	return fmt.Sprintf("%s as %s", t.Name, t.getAlias())
}

func (t *TableName) getAlias() string {
	if t.Alias == nil {
		return t.Name
	}
	return *t.Alias
}

func (t *TableName) GetTableType() string {
	return TableTypeTableName
}

type Subquery struct {
	Alias string
	Query Select
}

func (t *Subquery) ToSql() string {
	return fmt.Sprintf("%s as %s", t.Query.ToSql(), t.Alias)
}

func (t *Subquery) GetTableType() string {
	return TableTypeSubquery
}

type CrossJoinTable struct {
	left  Table
	right Table
}

func (t *CrossJoinTable) ToSql() string {
	return fmt.Sprintf("%s, %s", t.left.ToSql(), t.right.ToSql())
}

func (t *CrossJoinTable) GetTableType() string {
	return TableTypeCrossJoinTable
}

type JoinTable struct {
	JoinType string
	CrossJoinTable
	On Expression
}

func (t *JoinTable) ToSql() string {
	return fmt.Sprintf("%s %s %s ON %s", t.left.ToSql(), t.JoinType, t.right.ToSql(), t.On.ToSql())
}

func (t *JoinTable) GetTableType() string {
	return TableTypeJoinTable
}

const (
	StatementTypeSelect = "SELECT"
)

type Statement interface {
	ToSql() string
	GetStatementType() string
}

type Select struct {
	ReturnColumn []IColumn
	FromTable    Table
	Condition    *Expression
	GroupBy      *GroupByClause
	Having       *HavingClause
	OrderBy      *OrderByClause
	Limit        *Expression
	Offset       *Expression
	Union        []Select
}

func (s *Select) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("SELECT")
	for index, column := range s.ReturnColumn {
		sql.WriteString(column.ToSql())
		if index != len(s.ReturnColumn)-1 {
			sql.WriteString(", ")
		}
	}
	sql.WriteString(" FROM ")
	sql.WriteString(s.FromTable.ToSql())
	if s.Condition != nil {
		sql.WriteString(" WHERE ")
		sql.WriteString((*s.Condition).ToSql())
	}
	if s.GroupBy != nil {
		sql.WriteString(fmt.Sprintf(" %s", s.GroupBy.ToSql()))
	}
	if s.Having != nil {
		sql.WriteString(fmt.Sprintf(" %s", s.Having.ToSql()))
	}
	if s.OrderBy != nil {
		sql.WriteString(fmt.Sprintf(" %s", s.OrderBy.ToSql()))
	}
	if s.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %s", (*s.Limit).ToSql()))
	}
	if s.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %s", (*s.Offset).ToSql()))
	}
	if len(s.Union) > 0 {
		for _, union := range s.Union {
			sql.WriteString(" UNION ")
			sql.WriteString(union.ToSql())
		}
	}
	return sql.String()
}

func (s *Select) GetStatementType() string {
	return StatementTypeSelect
}

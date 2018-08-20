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
	Expr Expression
	Column
}

func (c *ExprColumn) ToSql() string {
	return fmt.Sprintf("%s as %s", c.Expr.ToSql(), c.getAlias())
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
	ExprTypeValueExpr     = "ValueExpr"
	ExprTypeColumnExpr    = "ColumnExpr"
	ExprTypeUnaryExpr     = "UnaryExpr"
	ExprTypeBinaryExpr    = "BinaryExpr"
	ExprTypeNotExpr       = "NotExpr"
	ExprTypeAndExpr       = "AndExpr"
	ExprTypeOrExpr        = "OrExpr"
	ExprTypeBetweenExpr   = "BetweenExpr"
	ExprTypeInExpr        = "InExpr"
	ExprTypeExistsExpr    = "ExistsExpr"
	ExprTypeValueListExpr = "ValueListExpr"
	ExprTypeSubqueryExpr  = "SubqueryExpr"
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
	Operator string
	Expr     Expression
}

func (e *UnaryExpr) ToSql() string {
	return fmt.Sprintf("%s%s", e.Operator, e.Expr.ToSql())
}

func (e *UnaryExpr) GetExprType() string {
	return ExprTypeUnaryExpr
}

type NotExpr struct {
	Expr Expression
}

func (e *NotExpr) ToSql() string {
	return "NOT " + e.Expr.ToSql()
}

func (e *NotExpr) GetExprType() string {
	return ExprTypeNotExpr
}

//二元运算，如四则运算
type BinaryExpr struct {
	Operator string
	Left     Expression
	Right    Expression
}

func (e *BinaryExpr) ToSql() string {
	return fmt.Sprintf("%s %s %s", e.Left.ToSql(), e.Operator, e.Right.ToSql())
}

func (e *BinaryExpr) GetExprType() string {
	return ExprTypeBinaryExpr
}

type AndExpr struct {
	Left  Expression
	Right Expression
}

func (e *AndExpr) ToSql() string {
	sql := strings.Builder{}
	if e.Left.GetExprType() == ExprTypeAndExpr || e.Left.GetExprType() == ExprTypeOrExpr {
		sql.WriteString(fmt.Sprintf("(%s)", e.Left.ToSql()))
	} else {
		sql.WriteString(e.Left.ToSql())
	}
	sql.WriteString(" AND ")
	if e.Right.GetExprType() == ExprTypeAndExpr || e.Right.GetExprType() == ExprTypeOrExpr {
		sql.WriteString(fmt.Sprintf("(%s)", e.Right.ToSql()))
	} else {
		sql.WriteString(e.Right.ToSql())
	}
	return sql.String()
}

func (e *AndExpr) GetExprType() string {
	return ExprTypeAndExpr
}

type OrExpr struct {
	AndExpr
}

func (e *OrExpr) ToSql() string {
	sql := strings.Builder{}
	if e.Left.GetExprType() == ExprTypeAndExpr || e.Left.GetExprType() == ExprTypeOrExpr {
		sql.WriteString(fmt.Sprintf("(%s)", e.Left.ToSql()))
	} else {
		sql.WriteString(e.Left.ToSql())
	}
	sql.WriteString(" OR ")
	if e.Right.GetExprType() == ExprTypeAndExpr || e.Right.GetExprType() == ExprTypeOrExpr {
		sql.WriteString(fmt.Sprintf("(%s)", e.Right.ToSql()))
	} else {
		sql.WriteString(e.Right.ToSql())
	}
	return sql.String()
}

func (e *OrExpr) GetExprType() string {
	return ExprTypeOrExpr
}

type BetweenExpr struct {
	Not   bool
	Expr  Expression
	Left  Expression
	Right Expression
}

func (e *BetweenExpr) ToSql() string {
	if e.Not {
		return fmt.Sprintf("%s NOT BETWEEN %s AND %s", e.Expr.ToSql(), e.Left.ToSql(), e.Right.ToSql())
	}
	return fmt.Sprintf("%s BETWEEN %s AND %s", e.Expr.ToSql(), e.Left.ToSql(), e.Right.ToSql())
}

func (e *BetweenExpr) GetExprType() string {
	return ExprTypeBetweenExpr
}

type InExpr struct {
	Not  bool
	Expr Expression
}

func (e *InExpr) ToSql() string {
	if e.Not {
		return fmt.Sprintf("%s NOT IN", e.Expr.ToSql())
	}
	return fmt.Sprintf("%s IN", e.Expr.ToSql())
}

func (e *InExpr) GetExprType() string {
	return ExprTypeInExpr
}

type ValueListExpr struct {
	Values []Expression
}

func (e *ValueListExpr) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("(")
	for index, value := range e.Values {
		sql.WriteString(value.ToSql())
		if index != len(e.Values)-1 {
			sql.WriteString(", ")
		}
	}
	sql.WriteString(")")
	return sql.String()
}

func (e *ValueListExpr) GetExprType() string {
	return ExprTypeValueListExpr
}

type ExistsExpr struct {
	Not  bool
	Expr Expression
}

func (e *ExistsExpr) ToSql() string {
	if e.Not {
		return fmt.Sprintf("%s NOT EXISTS", e.Expr.ToSql())
	}
	return fmt.Sprintf("%s EXISTS", e.Expr.ToSql())
}

func (e *ExistsExpr) GetExprType() string {
	return ExprTypeExistsExpr
}

type SubqueryExpr struct {
	Expr Select
}

func (e *SubqueryExpr) ToSql() string {
	return fmt.Sprintf("(%s)", e.Expr.ToSql())
}

func (e *SubqueryExpr) GetExprType() string {
	return ExprTypeSubqueryExpr
}

type GroupByClause struct {
	Columns []Expression
}

func (g *GroupByClause) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("GROUP BY ")
	for index, column := range g.Columns {
		sql.WriteString(fmt.Sprintf(" %s", column.ToSql()))
		if index != len(g.Columns)-1 {
			sql.WriteString(",")
		}
	}
	return sql.String()
}

type HavingClause struct {
	Expr Expression
}

func (h *HavingClause) ToSql() string {
	return fmt.Sprintf("HAVING %s", h.Expr.ToSql())
}

type OrderByClause struct {
	Columns []OrderByColumn
}

type OrderByColumn struct {
	Expr Expression
	Asc  bool
}

func (o *OrderByClause) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("ORDER BY ")
	for index, column := range o.Columns {
		if column.Asc {
			sql.WriteString(fmt.Sprintf("%s ASC", column.Expr.ToSql()))
		} else {
			sql.WriteString(fmt.Sprintf("%s DESC", column.Expr.ToSql()))
		}
		if index != len(o.Columns)-1 {
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
	return fmt.Sprintf("(%s) as %s", t.Query.ToSql(), t.Alias)
}

func (t *Subquery) GetTableType() string {
	return TableTypeSubquery
}

type CrossJoinTable struct {
	Left  Table
	Right Table
}

func (t *CrossJoinTable) ToSql() string {
	return fmt.Sprintf("%s, %s", t.Left.ToSql(), t.Right.ToSql())
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
	return fmt.Sprintf("%s %s %s ON %s", t.Left.ToSql(), t.JoinType, t.Right.ToSql(), t.On.ToSql())
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
	Condition    Expression
	GroupBy      *GroupByClause
	Having       *HavingClause
	OrderBy      *OrderByClause
	Limit        Expression
	Offset       Expression
	Union        []Select
}

func (s *Select) ToSql() string {
	sql := strings.Builder{}
	sql.WriteString("SELECT ")
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
		sql.WriteString(s.Condition.ToSql())
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
		sql.WriteString(fmt.Sprintf(" LIMIT %s", s.Limit.ToSql()))
	}
	if s.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %s", s.Offset.ToSql()))
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

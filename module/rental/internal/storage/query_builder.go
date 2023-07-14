package storage

import (
	"strconv"
	"strings"
)

type QueryBuilder struct {
	queryType   string
	targetTable string
	joins       []string
	columns     []string
	conditions  []string
	limit       *int
	offset      *int
	orderBy     *string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) Select() *QueryBuilder {
	qb.queryType = "SELECT "
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.targetTable = table
	return qb
}

func (qb *QueryBuilder) Join(joins ...string) *QueryBuilder {
	qb.joins = append(qb.joins, joins...)
	return qb
}

func (qb *QueryBuilder) Columns(columns ...string) *QueryBuilder {
	qb.columns = append(qb.columns, columns...)
	return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.conditions = append(qb.conditions, condition)
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = &limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = &offset
	return qb
}

func (qb *QueryBuilder) OrderBy(orderBy string) *QueryBuilder {
	qb.orderBy = &orderBy
	return qb
}

func (qb *QueryBuilder) String() string {
	var sb strings.Builder

	sb.WriteString(qb.queryType)
	sb.WriteString(strings.Join(qb.columns, ","))
	sb.WriteString(" FROM ")
	sb.WriteString(qb.targetTable)

	if len(qb.joins) > 0 {
		for _, join := range qb.joins {
			sb.WriteString(" JOIN ")
			sb.WriteString(join)
		}
	}

	if len(qb.conditions) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.conditions, " AND "))
	}

	if qb.limit != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(strconv.Itoa(*qb.limit))
	}

	if qb.offset != nil {
		sb.WriteString(" OFFSET ")
		sb.WriteString(strconv.Itoa(*qb.offset))
	}

	if qb.orderBy != nil {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(*qb.orderBy)
	}

	return sb.String()
}

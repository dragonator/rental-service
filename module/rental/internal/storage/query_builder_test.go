package storage_test

import (
	"strings"
	"testing"

	"github.com/dragonator/rental-service/module/rental/internal/storage"
)

func TestQueryBuilder(t *testing.T) {
	tests := []struct {
		name             string
		queryBuilderFunc func(*storage.QueryBuilder) *storage.QueryBuilder
		expectedQuery    string
	}{
		{
			name: "Basic select",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users")
			},
			expectedQuery: "SELECT * FROM users",
		},
		{
			name: "Select multiple columns #1",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("id", "name", "age").
					From("users")
			},
			expectedQuery: "SELECT id, name, age FROM users",
		},
		{
			name: "Select multiple columns #2",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("id").
					Columns("name", "age").
					From("users")
			},
			expectedQuery: "SELECT id, name, age FROM users",
		},
		{
			name: "Join tables",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users").
					Join("orders ON users.id = orders.user_id").
					Join("payments ON users.id = payments.user_id")
			},
			expectedQuery: "SELECT * FROM users JOIN orders ON users.id = orders.user_id JOIN payments ON users.id = payments.user_id",
		},
		{
			name: "Where",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users").
					Where("age > 18").
					Where("country = 'USA'")
			},
			expectedQuery: "SELECT * FROM users WHERE age > 18 AND country = 'USA'",
		},
		{
			name: "Limit",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users").
					Limit(10)
			},
			expectedQuery: "SELECT * FROM users LIMIT 10",
		},
		{
			name: "Offset",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users").
					Offset(20)
			},
			expectedQuery: "SELECT * FROM users OFFSET 20",
		},
		{
			name: "OrderBy",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("*").
					From("users").
					OrderBy("name ASC")
			},
			expectedQuery: "SELECT * FROM users ORDER BY name ASC",
		},
		{
			name: "Complex",
			queryBuilderFunc: func(qb *storage.QueryBuilder) *storage.QueryBuilder {
				return qb.Select().
					Columns("users.name", "orders.order_id", "payments.amount").
					From("users").
					Join("orders ON users.id = orders.user_id").
					Join("payments ON users.id = payments.user_id").
					Where("users.age > 18").
					OrderBy("users.name ASC").
					Limit(10).
					Offset(20)
			},
			expectedQuery: "SELECT users.name, orders.order_id, payments.amount FROM users JOIN orders ON users.id = orders.user_id JOIN payments ON users.id = payments.user_id WHERE users.age > 18 ORDER BY users.name ASC LIMIT 10 OFFSET 20",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			qb := storage.NewQueryBuilder()
			qb = test.queryBuilderFunc(qb)
			query := qb.String()

			if strings.TrimSpace(query) != test.expectedQuery {
				t.Errorf("\nExpected: %s\ngot: %s\n", test.expectedQuery, query)
			}
		})
	}
}

package data

import "github.com/Masterminds/squirrel"

// CreateSqlBuilder will create a squirrel SQL builder with preset for postgres
func CreateSqlBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

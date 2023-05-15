package data

import "github.com/lib/pq"

const (
	UniqueViolationError     = pq.ErrorCode("23505") // 'unique_violation'
	NotNullViolationError    = pq.ErrorCode("23502") // 'not_null_violation'
	ForeignKeyViolationError = pq.ErrorCode("23503") // 'foreign_key_violation'
)

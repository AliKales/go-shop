package utils

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsNotUniqueColumn(err error, indexName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, indexName) {
		return true
	}
	return false
}

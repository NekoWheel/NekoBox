package dbutil

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

func IsUniqueViolation(err error, constraint string) bool {
	if err == nil {
		return false
	}

	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	if sqlErr.Number == 1062 {
		// MySQL error code 1062 is ER_DUP_ENTRY, indicating a duplicate entry error
		// Extract the conflicting index name from the error message
		// The error message format is like this:
		// "Duplicate entry '{value}' for key '{index_name}'"
		msg := sqlErr.Message
		i := strings.Index(msg, "for key '")
		if i == -1 {
			return false
		}
		j := strings.Index(msg[i+len("for key '"):], "'")
		if j == -1 {
			return false
		}
		indexName := msg[i+len("for key '") : i+len("for key '")+j]
		return indexName == constraint
	}

	return false
}

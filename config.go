package sqldimel

import (
	"database/sql"
)

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
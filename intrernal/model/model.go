package model

import "database/sql"

type Url struct {
	ID            int64
	Alias         string
	URL           string
	RedirectCount sql.NullInt64
}

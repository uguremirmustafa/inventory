package utils

import (
	"database/sql"
	"time"
)

func GetNilTime(t *sql.NullTime) *time.Time {
	var timeData *time.Time
	if t.Valid {
		timeData = &t.Time
	} else {
		timeData = nil
	}
	return timeData
}

func GetNilString(s *sql.NullString) *string {
	var str *string
	if s.Valid {
		str = &s.String
	} else {
		str = nil
	}
	return str
}

func GetNilInt64(i64 *sql.NullInt64) *int64 {
	var i *int64
	if i64.Valid {
		i = &i64.Int64
	} else {
		i = nil
	}
	return i
}

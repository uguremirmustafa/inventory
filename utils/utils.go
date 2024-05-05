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

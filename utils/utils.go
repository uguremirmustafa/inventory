package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

func GetPathID(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	idInt64, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("Cannot convert id path value: ", slog.String("id", id))
		return 0, fmt.Errorf("internal server error")
	}
	return int64(idInt64), nil
}

func GenerateToken() (string, error) {
	// Generate 32 bytes of random data
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a URL-safe string
	tokenString := base64.URLEncoding.EncodeToString(token)
	return tokenString, nil
}

func IsFalsy(s string) bool {
	return s == ""
}

package api

import (
	"context"
	"net/http"

	"github.com/uguremirmustafa/inventory/db"
)

func handleGreet(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := q.GetUser(context.Background(), 1)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write([]byte(user.Name))
	})
}

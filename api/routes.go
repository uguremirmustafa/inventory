package api

import (
	"database/sql"
	"net/http"

	"github.com/justinas/alice"
	"github.com/uguremirmustafa/inventory/db"
)

func addRoutes(mux *http.ServeMux, q *db.Queries, db *sql.DB) {
	chain := alice.New(logMiddleware)
	authChain := alice.New(logMiddleware, authMiddleware())

	authService := NewAuthService(q, db)
	mux.Handle("POST /v1/auth/login", chain.Then(Make(authService.HandleLogin)))
	mux.Handle("GET /v1/auth/logout", authChain.Then(Make(authService.HandleLogout)))
	mux.Handle("GET /v1/me", authChain.Then(handleMe(q)))

	itemTypeService := NewItemTypeService(q)
	mux.Handle("GET /v1/item-type", authChain.Then(Make(itemTypeService.HandleListItemType)))
	mux.Handle("POST /v1/item-type", authChain.Then(Make(itemTypeService.HandleCreateItemType)))

	mux.Handle("GET /v1/manufacturer", authChain.Then(handleListManufacturer(q)))

	locationService := NewLocationService(q, db)
	mux.Handle("GET /v1/location", authChain.Then(Make(locationService.HandleListUserLocations)))
	mux.Handle("GET /v1/location/{id}", authChain.Then(Make(locationService.HandleGetLocation)))

	itemService := NewItemService(q, db)
	mux.Handle("GET /v1/item", authChain.Then(Make(itemService.HandleListUserItem)))
	mux.Handle("POST /v1/item", authChain.Then(Make(itemService.HandleInsertUserItem)))

	uploadService := NewUploadService(q, db)
	mux.Handle("POST /v1/upload-images", authChain.Then(Make(uploadService.HandleUploadImages)))
}

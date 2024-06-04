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

	// TODO: make this path relative??
	uploadDir := http.Dir("/home/anomy/Dev/personal/homventory/backend/uploads")
	// uploadDir := http.Dir("../uploads")
	fileServer := http.FileServer(uploadDir)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fileServer))

	authService := NewAuthService(q, db)
	mux.Handle("GET /v1/auth/login-google", chain.Then(Make(authService.HandleLoginWithGoogle)))
	mux.Handle("GET /v1/auth/google-callback", chain.Then(Make(authService.HandleGoogleCallback)))
	mux.Handle("GET /v1/auth/logout", authChain.Then(Make(authService.HandleLogout)))
	mux.Handle("GET /v1/me", authChain.Then(handleMe(q)))

	groupService := NewGroupsService(q, db)
	mux.Handle("GET /v1/groups", authChain.Then(Make(groupService.HandleGetGroupsOfUser)))
	mux.Handle("PUT /v1/user-group", authChain.Then(Make(groupService.HandleUpdateActiveGroupOfUser)))
	mux.Handle("GET /v1/group-members", authChain.Then(Make(groupService.HandleGetGroupMembers)))

	itemTypeService := NewItemTypeService(q)
	mux.Handle("GET /v1/item-type", authChain.Then(Make(itemTypeService.HandleListItemTypes)))
	mux.Handle("POST /v1/item-type", authChain.Then(Make(itemTypeService.HandleCreateItemType)))

	mux.Handle("GET /v1/manufacturer", authChain.Then(handleListManufacturer(q)))

	locationService := NewLocationService(q, db)
	mux.Handle("GET /v1/location", authChain.Then(Make(locationService.HandleListUserLocations)))
	mux.Handle("GET /v1/location/{id}", authChain.Then(Make(locationService.HandleGetLocation)))
	mux.Handle("POST /v1/location", authChain.Then(Make(locationService.HandleInsertUserLocation)))
	mux.Handle("PUT /v1/location/{id}", authChain.Then(Make(locationService.HandleUpdateUserLocation)))
	mux.Handle("DELETE /v1/location/{id}", authChain.Then(Make(locationService.HandleDeleteUserLocation)))

	itemService := NewItemService(q, db)
	mux.Handle("GET /v1/item", authChain.Then(Make(itemService.HandleListItems)))
	mux.Handle("POST /v1/item", authChain.Then(Make(itemService.HandleInsertUserItem)))

	uploadService := NewUploadService(q, db)
	mux.Handle("POST /v1/upload-images", authChain.Then(Make(uploadService.HandleUploadImages)))

	invitationService := NewInvitationService(q, db)
	mux.Handle("POST /v1/invite-user", authChain.Then(Make(invitationService.HandleCreateInvitation)))

}

package api

import (
	"database/sql"
	"net/http"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/utils"
)

type InvitationService struct {
	q  *db.Queries
	db *sql.DB
}

func NewInvitationService(q *db.Queries, db *sql.DB) *InvitationService {
	return &InvitationService{
		q:  q,
		db: db,
	}
}

func (s *InvitationService) HandleCreateInvitation(w http.ResponseWriter, r *http.Request) error {
	groupID := getUserActiveGroupID(w, r)
	userID := getUserID(w, r)
	var reqBody struct {
		Email string `json:"email"`
	}
	err := decode(r, &reqBody)
	if err != nil {
		return InvalidJSON()
	}
	token, err := utils.GenerateToken()
	if err != nil {
		return err
	}

	invitation, err := s.q.CreateInvitation(r.Context(), db.CreateInvitationParams{
		Email:     reqBody.Email,
		Token:     token,
		GroupID:   groupID,
		InvitorID: userID,
	})
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, invitation)
}

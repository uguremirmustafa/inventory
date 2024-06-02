package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/internal/config"
	"github.com/uguremirmustafa/inventory/utils"
)

type GroupsService struct {
	q  *db.Queries
	db *sql.DB
}

func NewGroupsService(q *db.Queries, db *sql.DB) *GroupsService {
	return &GroupsService{
		q:  q,
		db: db,
	}
}

type Group struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

// HandleGetGroupsOfUser returns a list of Group
func (s *GroupsService) HandleGetGroupsOfUser(w http.ResponseWriter, r *http.Request) error {
	userID := getUserID(w, r)
	groups, err := s.q.GetGroupsOfUser(r.Context(), userID)
	if err != nil {
		return NotFound()
	}

	var list []Group
	for _, item := range groups {
		list = append(list, *getUserGroupJson(item))
	}
	return writeJson(w, http.StatusOK, list)
}

// HandleUpdateActiveGroupOfUser updates the current user's activeGroupID(switch family)
func (s *GroupsService) HandleUpdateActiveGroupOfUser(w http.ResponseWriter, r *http.Request) error {
	c := config.GetConfig()
	userID := getUserID(w, r)
	var reqBody struct {
		NewGroupID int64 `json:"newGroupID"`
	}
	err := decode(r, &reqBody)
	if err != nil {
		return InvalidJSON()
	}
	updatedUser, err := s.q.UpdateUserActiveGroupID(r.Context(), db.UpdateUserActiveGroupIDParams{
		ID:            userID,
		ActiveGroupID: sql.NullInt64{Int64: reqBody.NewGroupID, Valid: true},
	})
	if err != nil {
		return FailedUpsert()
	}
	jwtToken, err := createJWTToken(
		int(updatedUser.ID),
		updatedUser.Email,
		updatedUser.ActiveGroupID.Int64,
		[]byte(c.JwtSecret))
	if err != nil {
		return err
	}
	setAuthCookie(w, jwtToken, updatedUser)

	return writeJson(w, http.StatusOK, getUserJson(updatedUser))
}

type GroupMemberItem struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Avatar    *string `json:"avatar"`
	GroupID   int64   `json:"group_id"`
	GroupName string  `json:"group_name"`
}

func (s *GroupsService) HandleGetGroupMembers(w http.ResponseWriter, r *http.Request) error {
	groupID := getUserActiveGroupID(w, r)
	groupMemberItems, err := s.q.GetMembersOfGroup(r.Context(), groupID)
	if err != nil {
		return NotFound()
	}
	var list []GroupMemberItem
	for _, item := range groupMemberItems {
		list = append(list, *getMemberOfGroupJson(item))
	}
	return writeJson(w, http.StatusOK, list)
}

func getUserGroupJson(l db.GetGroupsOfUserRow) *Group {
	return &Group{
		ID:          l.GroupID,
		Name:        l.GroupName,
		Description: utils.GetNilString(&l.GroupDesc),
		CreatedAt:   utils.GetNilTime(&l.CreatedAt),
		UpdatedAt:   utils.GetNilTime(&l.UpdatedAt),
	}
}

func getMemberOfGroupJson(l db.GetMembersOfGroupRow) *GroupMemberItem {
	return &GroupMemberItem{
		ID:        l.ID,
		Name:      l.Name,
		Avatar:    utils.GetNilString(&l.Avatar),
		GroupID:   l.GroupID,
		GroupName: l.GroupName,
		Email:     l.Email,
	}
}

package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/uguremirmustafa/inventory/db"
	"github.com/uguremirmustafa/inventory/internal/config"
	"github.com/uguremirmustafa/inventory/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	q  *db.Queries
	db *sql.DB
}

func NewAuthService(q *db.Queries, db *sql.DB) *AuthService {
	return &AuthService{
		q:  q,
		db: db,
	}
}

type UserInfo struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type CtxUserID string
type CtxUserActiveGroupID string

const (
	ctxUserID            CtxUserID            = "userID"
	ctxUserActiveGroupID CtxUserActiveGroupID = "activeGroupID"
)

func (s *AuthService) upsertUserWithGroup(u *UserInfoResponse, ctx context.Context) (*db.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// transactional queries
	qtx := s.q.WithTx(tx)
	var groupID int64
	user, err := qtx.GetUserByEmail(ctx, u.Email)
	if err != nil {
		userParams := db.UpsertUserParams{
			Name:          u.Name,
			Email:         u.Email,
			Avatar:        sql.NullString{String: u.Picture, Valid: true},
			ActiveGroupID: sql.NullInt64{Valid: false},
		}
		user, err = qtx.UpsertUser(ctx, userParams)
		if err != nil {
			return nil, err
		}
	}
	groupID = user.ActiveGroupID.Int64

	if !user.ActiveGroupID.Valid || user.ActiveGroupID.Int64 <= 0 {
		slog.Debug(
			"activeGroupID is not valid or 0(zero)",
			slog.Int64("activeGroupID", user.ActiveGroupID.Int64),
		)
		// create group
		groupParams := db.CreateGroupParams{
			Name:         fmt.Sprintf("%s's Family", u.Name),
			GroupOwnerID: user.ID,
		}
		dbGroup, err := qtx.CreateGroup(ctx, groupParams)
		if err != nil {
			slog.Error(
				"error while creating the user's default group",
				slog.Int64("userID", user.ID),
			)
			return nil, err
		}

		// connect user and group using user_groups table
		userGroupParams := db.ConnectUserAndGroupParams{
			UserID:  user.ID,
			GroupID: dbGroup.ID,
		}
		err = qtx.ConnectUserAndGroup(ctx, userGroupParams)
		if err != nil {
			slog.Error(
				"error while connecting user to its group",
				slog.Int64("userID", user.ID),
				slog.Int64("groupID", dbGroup.ID),
			)
			return nil, err
		}
		user, err = qtx.UpdateUserActiveGroupID(ctx, db.UpdateUserActiveGroupIDParams{
			ID:            user.ID,
			ActiveGroupID: sql.NullInt64{Valid: true, Int64: dbGroup.ID},
		})
		if err != nil {
			slog.Error(
				"sth went wrong updating user's ActiveGroupID",
				slog.Int64("userID", user.ID),
				slog.Int64("groupID", dbGroup.ID),
			)
		}
		groupID = user.ActiveGroupID.Int64
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("transaction error")
		return nil, err
	}
	return &db.User{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Avatar:        user.Avatar,
		ActiveGroupID: sql.NullInt64{Valid: true, Int64: groupID},
	}, nil
}

func getoauthConfGoogle() *oauth2.Config {
	c := config.GetConfig()
	oauthConfGoogle := &oauth2.Config{
		ClientID:     c.GoogleClientID,
		ClientSecret: c.GoogleClientSecret,
		RedirectURL:  c.GoogleAuthRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
	return oauthConfGoogle
}

func (s *AuthService) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	c := config.GetConfig()
	oauthConfGoogle := getoauthConfGoogle()
	url := oauthConfGoogle.AuthCodeURL(c.GoogleOauthStateString)
	return writeJson(w, http.StatusOK, url)
}

func (s *AuthService) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) error {
	c := config.GetConfig()
	errorUrl := c.ClientAuthErrorCallback
	oauthConfGoogle := getoauthConfGoogle()
	state := r.FormValue("state")
	if state != c.GoogleOauthStateString {
		slog.Error("state string does not match", slog.String("stateString", state))
		http.Redirect(w, r, errorUrl, http.StatusTemporaryRedirect)
		return nil
	}

	code := r.FormValue("code")
	token, err := oauthConfGoogle.Exchange(context.Background(), code)
	if err != nil {
		slog.Error("Error exchanging code", slog.String("err", err.Error()))
		http.Redirect(w, r, errorUrl, http.StatusTemporaryRedirect)
		return err
	}

	// Use the access token to fetch user info
	userInfo, err := getUserInfo(token)
	if err != nil {
		slog.Error("Error getting user info from google", slog.String("err", err.Error()))
		http.Redirect(w, r, errorUrl, http.StatusTemporaryRedirect)
		return err
	}

	dbUser, err := s.upsertUserWithGroup(userInfo, r.Context())
	if err != nil {
		slog.Error("Error upserting user with group", slog.String("err", err.Error()))
		http.Redirect(w, r, errorUrl, http.StatusTemporaryRedirect)
		return err
	}

	slog.Debug("upsertUserWithGroup result", slog.Any("dbUser", dbUser))

	jwtToken, err := createJWTToken(
		int(dbUser.ID),
		dbUser.Email,
		dbUser.ActiveGroupID.Int64,
		[]byte(c.JwtSecret))
	if err != nil {
		return err
	}
	setAuthCookie(w, jwtToken, *dbUser)

	http.Redirect(w, r, c.ClientProfilePage, http.StatusTemporaryRedirect)
	return nil

}

func (s *AuthService) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	c := config.GetConfig()

	userID := getUserID(w, r)
	http.SetCookie(w, &http.Cookie{
		Name:     c.JwtCookieKey,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	slog.Info("user logged out", slog.Int64("userID", userID))
	writeJson(w, http.StatusOK, "logout successfull")
	return nil
}

func handleMe(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(w, r)

		user, err := q.GetUser(r.Context(), userID)
		if err != nil {
			slog.Error("user not found. Email: %s, ID: %v", user.Email, user.ID)
			writeJson(w, http.StatusNotFound, "user not found")
			return
		}

		writeJson(w, http.StatusOK, getUserJson(user))
	})
}

func authMiddleware() Middleware {
	c := config.GetConfig()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(Make(func(w http.ResponseWriter, r *http.Request) error {
			// Get JWT token from the cookie
			cookie, err := r.Cookie(c.JwtCookieKey)
			if err != nil {
				slog.Warn("no cookie found on request")
				return NotAuthorized()
			}

			// Validate JWT token
			tokenString := cookie.Value
			claims := &MyCustomClaims{}
			err = verifyToken(tokenString, []byte(c.JwtSecret), claims)
			if err != nil {
				slog.Warn("token verification failed.", slog.String("token", tokenString))
				return NotAuthorized()
			}

			// Check token expiry
			if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
				slog.Warn("token expired.", slog.String("token", tokenString))
				return NotAuthorized()
			}

			if claims.ActiveGroupID == 0 {
				slog.Error("no active group id", slog.Int64("activeGroupID", claims.ActiveGroupID))
			}
			// JWT token is valid, proceed with the next handler
			ctx := context.WithValue(r.Context(), ctxUserID, claims.Subject)
			ctx = context.WithValue(ctx, ctxUserActiveGroupID, claims.ActiveGroupID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}))
	}
}

type User struct {
	ID            int64   `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	Avatar        *string `json:"avatar"`
	ActiveGroupID *int64  `json:"activeGroupID"`
}

func getUserJson(l db.User) *User {
	return &User{
		ID:            l.ID,
		Name:          l.Name,
		Email:         l.Email,
		Avatar:        utils.GetNilString(&l.Avatar),
		ActiveGroupID: utils.GetNilInt64(&l.ActiveGroupID),
	}
}

func setAuthCookie(w http.ResponseWriter, token string, user db.User) {
	c := config.GetConfig()
	// Set JWT token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     c.JwtCookieKey,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(24 * time.Hour),
		Path:     "/",
	})
	slog.Info("created jwtToken", slog.String("token", token), slog.String("email", user.Email), slog.Int64("userID", user.ID))
}

func getUserID(w http.ResponseWriter, r *http.Request) int64 {
	ctx := r.Context()
	value := ctx.Value(ctxUserID).(string)
	userID, err := strconv.Atoi(value)
	if err != nil {
		slog.Error("Cannot convert userID: ", slog.Int("userID", userID))
		writeJson(w, http.StatusUnauthorized, "unauthorized")
	}
	return int64(userID)
}

func getUserActiveGroupID(w http.ResponseWriter, r *http.Request) int64 {
	ctx := r.Context()
	value := ctx.Value(ctxUserActiveGroupID).(int64)
	if value == 0 {
		slog.Error("Active group id is 0")
		writeJson(w, http.StatusUnauthorized, "unauthorized")
	}
	return value
}

type MyCustomClaims struct {
	Email         string `json:"email"`
	ActiveGroupID int64  `json:"activeGroupID"`
	jwt.RegisteredClaims
}

func createJWTToken(userID int, email string, activeGroupID int64, secret []byte) (string, error) {
	claims := MyCustomClaims{
		email,
		activeGroupID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func verifyToken(tokenString string, secret []byte, myClaims *MyCustomClaims) error {
	token, err := jwt.ParseWithClaims(tokenString, myClaims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		fmt.Printf("%+v\n", claims.ID)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}
	return nil
}

type UserInfoResponse struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}

// Create HTTP client with the access token
func getUserInfo(token *oauth2.Token) (*UserInfoResponse, error) {
	// Create HTTP client with the access token
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	response, err := httpClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	var userInfo UserInfoResponse
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

package api

import (
	"context"
	"database/sql"
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

const (
	ctxUserID CtxUserID = "userID"
)

func (s *AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	c := config.GetConfig()
	req := &UserInfo{}
	err := decode(r, req)
	if err != nil {
		return err
	}

	userParams := db.UpsertUserParams{
		Name:   req.Name,
		Email:  req.Email,
		Avatar: sql.NullString{String: req.Avatar, Valid: true},
	}

	user, err := s.q.UpsertUser(r.Context(), userParams)
	if err != nil {
		return err
	}

	// Create jwt token
	jwtToken, err := createJWTToken(int(user.ID), user.Email, []byte(c.JwtSecret))
	if err != nil {
		return err
	}
	setAuthCookie(w, jwtToken, user)

	return writeJson(w, http.StatusOK, getUserJson(user))
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

			// JWT token is valid, proceed with the next handler
			ctx := context.WithValue(r.Context(), ctxUserID, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
			return nil
		}))
	}
}

type User struct {
	ID     int64   `json:"id"`
	Email  string  `json:"email"`
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}

func getUserJson(l db.User) *User {
	return &User{
		ID:     l.ID,
		Name:   l.Name,
		Email:  l.Email,
		Avatar: utils.GetNilString(&l.Avatar),
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

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func createJWTToken(userID int, email string, secret []byte) (string, error) {
	claims := MyCustomClaims{
		email,
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

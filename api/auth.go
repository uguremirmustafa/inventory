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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfoResponse struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}

type CtxUserID string

const (
	ctxUserID CtxUserID = "userID"
)

func handleLoginGoogle() http.Handler {
	c := config.GetConfig()
	googleOauthConfig := getGoogleAuthConfig(c)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(c.GoogleOauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}
func handleCallbackGoogle(q *db.Queries) http.Handler {
	c := config.GetConfig()
	googleOauthConfig := getGoogleAuthConfig(c)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")

		if state != c.GoogleOauthStateString {
			slog.Error("Invalid oauth state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")

		token, err := googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			slog.Error("Error exchanging code", slog.String("error", err.Error()), slog.String("code", code))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Use the access token to fetch user info
		userInfo, err := getUserInfo(token)
		if err != nil {
			slog.Error("Error getting user info", slog.String("error", err.Error()))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Check if the user exists in the database
		user, err := q.GetUserByEmail(context.Background(), userInfo.Email)
		if err != nil {
			slog.Error("user not found, trying to insert")
			var u = db.CreateUserParams{
				Name:   userInfo.Name,
				Email:  userInfo.Email,
				Avatar: sql.NullString{String: userInfo.Picture, Valid: true},
			}
			user, err = q.CreateUser(context.Background(), u)
			if err != nil {
				slog.Error("Failed to create user")
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		}

		// Create jwt token
		jwtToken, err := createJWTToken(int(user.ID), user.Email, []byte(c.JwtSecret))
		if err != nil {
			slog.Error("Failed to create token")
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		// Set JWT token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     c.JwtCookieKey,
			Value:    jwtToken,
			HttpOnly: true,
			Expires:  time.Now().UTC().Add(24 * time.Hour),
		})
		slog.Info("created jwtToken", slog.String("token", jwtToken))

		http.Redirect(w, r, "/v1/me", http.StatusTemporaryRedirect)
	})
}

func handleMe(q *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(w, r)

		user, err := q.GetUser(r.Context(), userID)
		if err != nil {
			slog.Error("user not found. Email: %s, ID: %v", user.Email, user.ID)
			encode(w, http.StatusNotFound, "user not found")
			return
		}

		encode(w, http.StatusOK, getUserJson(user))
	})
}

func getUserID(w http.ResponseWriter, r *http.Request) int64 {
	ctx := r.Context()
	value := ctx.Value(ctxUserID).(string)
	userID, err := strconv.Atoi(value)
	if err != nil {
		slog.Error("Cannot convert userID: ", slog.Int("userID", userID))
		redirectToLogin(w, r)
	}
	return int64(userID)
}

// Create HTTP client with the access token
func getUserInfo(token *oauth2.Token) (*UserInfoResponse, error) {
	// Create HTTP client with the access token
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	// Make request to Google UserInfo endpoint
	response, err := httpClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		slog.Error("Failed to call to google for userinfo", slog.String("error", err.Error()))
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch user info")
	}

	// Decode JSON response into UserInfo struct
	var userInfo UserInfoResponse
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
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

func authMiddleware() Middleware {
	c := config.GetConfig()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT token from the cookie
			cookie, err := r.Cookie(c.JwtCookieKey)
			if err != nil {
				slog.Warn("no cookie found on request")
				redirectToLogin(w, r)
				return
			}

			// Validate JWT token
			tokenString := cookie.Value
			claims := &MyCustomClaims{}
			err = verifyToken(tokenString, []byte(c.JwtSecret), claims)
			if err != nil {
				slog.Warn("token verification failed.", slog.String("token", tokenString))
				redirectToLogin(w, r)
				return
			}

			// Check token expiry
			if time.Unix(claims.ExpiresAt.Unix(), 0).Before(time.Now()) {
				slog.Warn("token expired.", slog.String("token", tokenString))
				redirectToLogin(w, r)
				return
			}

			// JWT token is valid, proceed with the next handler
			ctx := context.WithValue(r.Context(), ctxUserID, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("redirecting to login")
	http.Redirect(w, r, "/v1/auth/login", http.StatusSeeOther)
}

type UserResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func getUserJson(dbUser db.User) UserResponse {
	responseData := UserResponse{
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Avatar: dbUser.Avatar.String,
	}
	return responseData
}

func getGoogleAuthConfig(c *config.Config) *oauth2.Config {
	googleOauthConfig := &oauth2.Config{
		ClientID:     c.GoogleClientID,
		ClientSecret: c.GoogleClientSecret,
		RedirectURL:  c.GoogleAuthRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return googleOauthConfig
}

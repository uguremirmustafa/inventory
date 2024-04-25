package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/uguremirmustafa/inventory/db"
	"golang.org/x/oauth2"
)

type UserInfoResponse struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Picture    string `json:"picture"`
}

func handleLoginGoogle(c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := c.googleOauthConfig.AuthCodeURL(c.oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}
func handleCallbackGoogle(q *db.Queries, c *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != c.oauthStateString {
			fmt.Println("Invalid oauth state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := c.googleOauthConfig.Exchange(context.Background(), code)
		if err != nil {
			fmt.Printf("Error exchanging code: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Use the access token to fetch user info
		userInfo, err := getUserInfo(token)
		if err != nil {
			fmt.Printf("Error getting user info: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// TODO:save/retrive user to/from db

		// Encode user info response to JSON
		userJson, err := json.Marshal(userInfo)
		if err != nil {
			fmt.Printf("Error encoding user info: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// Write JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	})
}

// Create HTTP client with the access token
func getUserInfo(token *oauth2.Token) (*UserInfoResponse, error) {
	// Create HTTP client with the access token
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	// Make request to Google UserInfo endpoint
	response, err := httpClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
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

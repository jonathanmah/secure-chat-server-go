package handlers

import (
	"chatapp/internal/auth"
	"chatapp/internal/config"
	"chatapp/internal/postgres"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// start authorization code flow
func RedirectOAuthHandler(w http.ResponseWriter, r *http.Request) {
	// generate random string
	state := generateRandomString()
	// add oauth random state on cookie
	http.SetCookie(w, config.NewOAuthStateCookie(state))
	url := config.App.Auth.OAuthConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, url, http.StatusFound)
}

// create a random string for oauth state
func generateRandomString() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err) // if this fails its doomed
	}
	return base64.URLEncoding.EncodeToString(b)
}

func PostOAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	// check state for CSRF, redirect must have same state it started with
	if !validateOAuthState(w, r) {
		return
	}
	// get authorization code from url
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}
	// get token from authorization server
	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	// get authenticated user info
	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	// get existing user ID or create one
	id, err := getID(userInfo.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// create jwt session token with user id
	sessionToken, err := auth.CreateSessionToken(id)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, config.NewSessionCookie(sessionToken))
	http.Redirect(w, r, config.App.Auth.PostOAuthRedirectURL, http.StatusFound)
}

// get ID from existing account or create OAuth account without password and return ID
func getID(email string) (string, error) {
	// check if account already been created with email
	emailExists, err := postgres.EmailExists(email)
	if err != nil {
		return "", errors.New("Failed checking database for email.")
	}
	var id string
	if emailExists {
		// if user already has an account, get their user id
		id, err = postgres.GetUserIdByEmail(email)
		if err != nil {
			return "", errors.New("Failed retrieving user ID by email from database.")
		}
	} else {
		// if no account exists for email, create new account without a password
		username, err := CreateUniqueUsername()
		if err != nil {
			return "", errors.New("Failed to create unique username.")
		}
		// get user id after postgres generates a new uuid for it
		id, err = postgres.CreatePasswordlessUser(email, username)
		if err != nil {
			return "", errors.New("Failed creating a new account in database.")
		}
	}
	return id, nil
}

// confirm oauth state is same after redirect
func validateOAuthState(w http.ResponseWriter, r *http.Request) bool {
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return false
	}
	return true
}

// trade authorization code for token
func exchangeCodeForToken(code string) (*oauth2.Token, error) {
	return config.App.Auth.OAuthConfig.Exchange(context.Background(), code)
}

type OAuthUserInfo struct {
	Sub   string `json:"sub"` // short for subject, whoever claim is for
	Email string `json:"email"`
	Name  string `json:"name"`
}

// get user info from google authority server
func fetchGoogleUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	client := config.App.Auth.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get(config.App.Auth.OAuthUserInfoURL) // #HARCODED
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

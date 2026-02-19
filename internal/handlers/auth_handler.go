package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"pg-management-system/internal/database"
	"pg-management-system/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

const oauthStateString = "random_state_string" // In production, use a more secure state management

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := getGoogleOauthConfig().AuthCodeURL(oauthStateString, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := getGoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	user, err := database.GetUserByEmail(googleUser.Email)
	if err != nil {
		http.Error(w, "Database search failed", http.StatusInternalServerError)
		return
	}

	if user == nil {
		user = &models.User{
			Email:    googleUser.Email,
			Name:     googleUser.Name,
			GoogleID: googleUser.ID,
		}
		if err := database.CreateUser(user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if user.GoogleID == "" {
		// Update existing user with Google ID
		if err := database.UpdateUserByGoogleID(googleUser.ID, googleUser.Name); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
		user.GoogleID = googleUser.ID
	}

	jwtToken, err := GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": jwtToken,
		"email": user.Email,
		"name":  user.Name,
	})
}

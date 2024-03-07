package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/handlers/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"forum/internal/domain"
	"forum/internal/helpers/cookies"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL    = "https://accounts.google.com/o/oauth2/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type googleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}

func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=profile email", googleAuthURL, h.googleConfig.ClientID, h.googleConfig.RedirectURL)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// Exchange authorization code for access token
	data := strings.NewReader(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, h.googleConfig.ClientID, h.googleConfig.ClientSecret, h.googleConfig.RedirectURL))
	resp, err := http.Post(googleTokenURL, "application/x-www-form-urlencoded", data)
	if err != nil {
		h.service.Log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.service.Log.Printf("Failed to read response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Extract access token from response
	accessToken := ExtractValueFromBody(body, "access_token")
	if accessToken == "" {
		h.service.Log.Printf("Failed to extract access token from response")
		http.Error(w, "Failed to extract access token", http.StatusInternalServerError)
		return
	}

	// Get user info using the access token
	userInfo, err := utils.GetUserInfo(accessToken, googleUserInfoURL)
	if err != nil {
		h.service.Log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	var googleUserInfo googleUserInfo

	err = json.Unmarshal(userInfo, &googleUserInfo)
	if err != nil {
		h.service.Log.Printf("Failed to unmarshal user info: %v", err)
		http.Error(w, "Failed to unmarshal user info", http.StatusInternalServerError)
		return
	}

	user, _ := h.service.UserService.GetUserByEmail(googleUserInfo.Email)
	if user == nil {
		userDTO := &domain.CreateUserDTO{
			Username: googleUserInfo.Name,
			Email:    googleUserInfo.Email,
			Password: googleUserInfo.Sub,
		}
		err = h.service.UserService.CreateUser(userDTO)
		if err != nil {
			h.service.Log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	userLogin := &domain.LoginUserDTO{
		Email:    googleUserInfo.Email,
		Password: googleUserInfo.Sub,
	}

	userID, err := h.service.UserService.LoginUser(userLogin)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.service.SessionService.CreateSession(userID)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.SetCookie(w, session.UUID, int(time.Until(session.ExpireAt).Seconds()))

	http.Redirect(w, r, "/", http.StatusFound)
}

func ExtractValueFromBody(body []byte, key string) string {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return ""
	}

	value, ok := response[key].(string)
	if !ok {
		return ""
	}

	return value
}

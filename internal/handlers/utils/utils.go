package utils

import (
	"fmt"
	"forum/internal/domain"
	"io/ioutil"
	"net/http"
)

type contextKey string

var ContextKeyUser = contextKey("user")

func GetUserFromContext(r *http.Request) *domain.User {
	user, ok := r.Context().Value(ContextKeyUser).(*domain.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserInfo(accessToken string, userInfoURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

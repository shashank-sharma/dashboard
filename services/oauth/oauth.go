package oauth

import (
	"encoding/json"
	"net/http"
)

type UserInfo struct {
	Email string `json:"email"`
}

func FetchUserInfo(client *http.Client) (*UserInfo, error) {
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var info UserInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

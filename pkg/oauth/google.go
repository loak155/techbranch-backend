package oauth

import (
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleManager struct {
	State        string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type userInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewGoogleManager(stat, clientID, clientSecret, redirectURL string) *GoogleManager {
	return &GoogleManager{
		State:        stat,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

func (g *GoogleManager) GenerateConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
}

func (g *GoogleManager) GetLoginURL() string {
	config := g.GenerateConfig()
	return config.AuthCodeURL(g.State, oauth2.AccessTypeOffline)
}

func (g *GoogleManager) CheckState(state string) bool {
	return state == g.State
}

func (g *GoogleManager) GetAccessToken(code string) (*oauth2.Token, error) {
	config := g.GenerateConfig()
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (g *GoogleManager) GetUserInfo(token *oauth2.Token) (userInfo, error) {
	userInfo := userInfo{}

	config := g.GenerateConfig()
	client := config.Client(oauth2.NoContext, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return userInfo, err
	}

	// userInfo := make(map[string]interface{})
	// err = json.NewDecoder(response.Body).Decode(&userInfo)
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

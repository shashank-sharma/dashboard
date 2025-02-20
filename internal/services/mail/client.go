package mail

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
)

type MailService struct {
	googleConfig *oauth2.Config
}

func NewMailService() *MailService {
	googleConfig := &oauth2.Config{
		ClientID:     os.Getenv("GCP_MAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GCP_MAIL_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GCP_MAIL_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes: []string{
			gmail.GmailReadonlyScope,
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}

	return &MailService{
		googleConfig: googleConfig,
	}
}

func (ms *MailService) GetAuthUrl() string {
	return ms.googleConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce, oauth2.SetAuthURLParam("prompt", "consent"), oauth2.SetAuthURLParam("redirect_uri", ms.googleConfig.RedirectURL))
}

func (ms *MailService) GetConfig() *oauth2.Config {
	return ms.googleConfig
}

func (ms *MailService) FetchClient(tokenId string) (*http.Client, error) {
	token, err := query.FindById[*models.Token](tokenId)
	if err != nil {
		return nil, err
	}

	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Time(),
	}

	logger.Debug.Println("Found oauthToken: ", oauthToken)

	return ms.googleConfig.Client(context.Background(), oauthToken), nil
}

func (ms *MailService) GetGmailService(client *http.Client) (*gmail.Service, error) {
	return gmail.NewService(context.Background(), option.WithHTTPClient(client))
}

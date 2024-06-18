package calendar

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func LoadCalendarConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GCP_CAL_CLIENT_ID"),
		ClientSecret: os.Getenv("GCP_CAL_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GCP_CAL_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarReadonlyScope, "https://www.googleapis.com/auth/userinfo.email"},
	}
}

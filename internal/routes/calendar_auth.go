package routes

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/models"
	"github.com/shashank-sharma/backend/internal/query"
	"github.com/shashank-sharma/backend/internal/services/calendar"
	"github.com/shashank-sharma/backend/internal/services/oauth"
	"github.com/shashank-sharma/backend/internal/util"
)

type CalendarTokenAPI struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

func CalendarAuthHandler(c echo.Context) error {
	googleConfig := calendar.LoadCalendarConfig()

	fmt.Println("GoogleConfig ", googleConfig)
	authUrl := googleConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", googleConfig.RedirectURL))
	return c.JSON(http.StatusOK, map[string]interface{}{"url": authUrl})
}

func CalendarAuthCallback(c echo.Context) error {
	pbToken := c.Request().Header.Get(echo.HeaderAuthorization)
	userId, err := util.GetUserId(pbToken)

	googleConfig := calendar.LoadCalendarConfig()

	calTokenData := &CalendarTokenAPI{}
	if err := c.Bind(calTokenData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	token, err := googleConfig.Exchange(context.Background(), calTokenData.Code)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Invalid token exchange"})
	}

	client := googleConfig.Client(context.Background(), token)
	userInfo, err := oauth.FetchUserInfo(client)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch userinfo"})
	}

	expiry := types.DateTime{}
	expiry.Scan(token.Expiry)

	calToken := &models.CalendarToken{
		User:         userId,
		Account:      userInfo.Email,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       expiry,
	}

	if err := query.UpsertRecord[*models.CalendarToken](calToken, map[string]interface{}{"account": userInfo.Email}); err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Error saving record"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Authenticated successfully"})
}

func TestRandomHandler(app *pocketbase.PocketBase, c echo.Context) error {
	logger.Debug.Println("Started track")
	user, err := query.FindById[*models.Users]("k0jhrgadiakg8xb")
	if err != nil {
		logger.Debug.Println("err: ", err)
	}
	logger.Debug.Println("user: ", user)
	logger.Debug.Println("Completed")

	logger.Debug.Println("User found:", user)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": user.Email})
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pocketbase/dbx"
	pocketbaseModel "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/shashank-sharma/backend/models"
	"github.com/shashank-sharma/backend/util"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
)

func streamMP3(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func streamMP3Handler(c echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "path parameter is required"})
	}

	reader, err := streamMP3(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error streaming mp3: %v", err)})
	}

	// Create a buffered reader with a buffer size of 4096 bytes.
	bufferedReader := bufio.NewReaderSize(reader, 4096)

	c.Response().Header().Set("Content-Type", "audio/mpeg")
	_, err = io.CopyBuffer(c.Response(), bufferedReader, nil)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	util.Init()
	app := util.GetApp()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		// prints "Hello!" every 2 minutes
		scheduler.MustAdd("hello", "*/6 * * * *", func() {
			threshold := time.Now().Add(-6 * time.Minute)
			records, err := app.Dao().FindRecordsByFilter(
				"devices", "is_online = {:is_online} && is_active = {:is_active} && updated <= {:threshold}",
				"",
				100,
				0,
				dbx.Params{"is_online": true},
				dbx.Params{"is_active": true},
				dbx.Params{"threshold": threshold})

			if err != nil {
				fmt.Println("No queries found")
			}
			for _, record := range records {
				record.Set("is_online", false)
				fmt.Println("Working for record =", record)
				if err := app.Dao().SaveRecord(record); err != nil {
					fmt.Println("Error saving cron: ", err)
				}
			}
		})

		scheduler.Start()
		// serves static files from the provided public dir (if exists)
		// e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/hello",
			Handler: func(c echo.Context) error {
				obj := map[string]interface{}{"message": "Hello world!"}
				return c.JSON(http.StatusOK, obj)
			},
			// Middlewares: []echo.MiddlewareFunc{
			// 	apis.RequireAdminOrUserAuth(),
			// },
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/token",
			Handler: func(c echo.Context) error {
				token := c.Request().Header.Get(echo.HeaderAuthorization)
				fmt.Println("token =", token)
				userId, err := util.GetUserId(token)
				if err != nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
				}
				record, err := app.Dao().FindFirstRecordByFilter(
					"dev_tokens", "user = {:user}",
					dbx.Params{"user": userId})

				// TODO: Create token if not found and return it
				if err != nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "No Dev token found"})
				}

				isActive := false
				devToken := record.Get("token")
				isActive = record.GetBool("is_active")
				if isActive {
					obj := map[string]interface{}{"token": devToken}
					fmt.Println("Successful token sent", obj)
					return c.JSON(http.StatusOK, obj)
				} else {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Dev token is disabled"})
				}

			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/track/create",
			Handler: func(c echo.Context) error {
				fmt.Println("Started track create")
				token := c.Request().Header.Get(echo.HeaderAuthorization)
				fmt.Println("token =", token)
				userId, err := util.GetUserId(token)
				if err != nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
				}

				data := &models.TrackDeviceAPI{}

				if err := c.Bind(data); err != nil || data.Name == "" {
					fmt.Println("Error in parsing =", err)
					return apis.NewBadRequestError("Failed to read request data", err)
				}
				record, _ := app.Dao().FindFirstRecordByFilter(
					"devices", "user = {:user} && name = {:name} && hostname = {:hostname} && os = {:os} && arch = {:arch}",
					dbx.Params{"user": userId},
					dbx.Params{"name": data.Name},
					dbx.Params{"hostname": data.HostName},
					dbx.Params{"os": data.Os},
					dbx.Params{"arch": data.Arch},
				)

				if record != nil {
					fmt.Println("returning id:", record.Get("id"))
					return c.JSON(http.StatusOK, map[string]interface{}{"id": record.Get("id")})
				} else {
					collection, err := app.Dao().FindCollectionByNameOrId("devices")
					if err != nil {
						return err
					}

					record := pocketbaseModel.NewRecord(collection)

					form := forms.NewRecordUpsert(app, record)
					formId := util.GenerateRandomId()
					form.LoadData(map[string]any{
						"id":        formId,
						"user":      userId,
						"name":      data.Name,
						"hostname":  data.HostName,
						"os":        data.Os,
						"arch":      data.Arch,
						"is_active": true,
						"is_online": true,
					})

					form.Validate()
					form.LoadRequest(c.Request(), "")
					if err := form.Submit(); err != nil {
						fmt.Println("Error saving file", err)
						return err
					}

					fmt.Println("returning formdata id:", formId)
					return c.JSON(http.StatusOK, map[string]interface{}{"id": formId})
				}

			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/track",
			Handler: func(c echo.Context) error {
				fmt.Println("Started track")
				data := &models.TrackDeviceUpdateAPI{}

				if err := c.Bind(data); err != nil || data.UserId == "" || data.Token == "" || data.ProductId == "" {
					fmt.Println("Error in parsing =", err)
					// TODO: Simply say unauthorized
					return apis.NewBadRequestError("Failed to read request data", err)
				}
				record, err := app.Dao().FindFirstRecordByFilter(
					"dev_tokens", "user = {:user} && token = {:token}",
					dbx.Params{"user": data.UserId},
					dbx.Params{"token": data.Token},
				)

				if err != nil || record == nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Not authorized"})
				}

				deviceRecord, err := app.Dao().FindFirstRecordByFilter(
					"devices", "user = {:user} && id = {:id}",
					dbx.Params{"user": data.UserId},
					dbx.Params{"id": data.ProductId},
				)

				if err != nil || deviceRecord == nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Authorized, but product not found"})
				}

				deviceRecord.Set("is_active", true)
				deviceRecord.Set("is_online", true)

				if err := app.Dao().SaveRecord(deviceRecord); err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err})
				}

				return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
			},
		})

		e.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/stream_mp3",
			Handler: streamMP3Handler,
		})

		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/sync/track-items",
			Handler: func(c echo.Context) error {
				token := c.Request().Header.Get(echo.HeaderAuthorization)
				fmt.Println("token =", token)
				userId, err := util.GetUserId(token)
				if err != nil {
					return c.JSON(http.StatusForbidden, map[string]interface{}{"message": "Failed to fetch id, token misconfigured"})
				}
				data := &models.TrackUploadAPI{}

				if err := c.Bind(data); err != nil || data.Source == "" {
					fmt.Println("Error in parsing =", err)
					return apis.NewBadRequestError("Failed to read request data", err)
				}
				//fileContent, err := c.FormFile("file")
				//if err != nil {
				//	return err
				//}
				collection, err := app.Dao().FindCollectionByNameOrId("track_upload")
				if err != nil {
					return err
				}

				record := pocketbaseModel.NewRecord(collection)

				form := forms.NewRecordUpsert(app, record)

				form.LoadData(map[string]any{
					"id":     util.GenerateRandomId(),
					"user":   userId,
					"status": "CREATED",
				})

				form.LoadRequest(c.Request(), "")
				fmt.Println("Checking form id =", form.Id, "record=", record.Id)
				if err := form.Submit(); err != nil {
					fmt.Println("Error saving file", err)
					return err
				}

				formData := form.Data()
				trackUpload := &models.TrackUpload{
					User:   userId,
					Source: formData["source"].(string),
					File:   formData["file"].(string),
					Synced: formData["synced"].(bool),
				}

				trackUpload.BaseModel.Id = form.Id

				// load the entire request

				go util.SyncTrackUpload(trackUpload, data.ForceCheck)
				return c.JSON(http.StatusOK, map[string]interface{}{"message": "Task scheduled", "track_upload": trackUpload})
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireRecordAuth(),
				apis.ActivityLogger(app),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

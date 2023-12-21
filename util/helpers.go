package util

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	pbModels "github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/shashank-sharma/backend/models"
)

type OperationCount struct {
	CreateCount int64 `json:"create_count"`
	SkipCount   int64 `json:"skip_count"`
	ForceCheck  bool  `json:"force_check"`
}

func GenerateRandomId() string {
	return security.RandomStringWithAlphabet(pbModels.DefaultIdLength, pbModels.DefaultIdAlphabet)
}

func GetUserId(tokenString string) (string, error) {
	// Split the token into header, payload, and signature
	parts := strings.Split(tokenString, ".")

	// Decode the payload (no signature verification)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("Error decoding payload:", err)
		return "", err
	}

	// Unmarshal the payload
	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		fmt.Println("Error unmarshalling payload:", err)
		return "", err
	}

	// Access the claims
	fmt.Println("User ID:", claims["id"])
	fmt.Println("Username:", claims["username"])
	return claims["id"].(string), nil
}

func SyncTrackUpload(trackUpload *models.TrackUpload, forceCheck bool) {
	trackUpload.Status = "IN-PROGRESS"
	trackUpload.MarkAsNotNew()
	if err := app.Dao().Save(trackUpload); err != nil {
		fmt.Println("Failed updating record")
		return
	}
	defer func() {
		if err := app.Dao().Save(trackUpload); err != nil {
			fmt.Println("Failed updating record")
		}
	}()

	opCount, err := insertFromFile(trackUpload, forceCheck)
	if err != nil {
		fmt.Println("Something went wrong while insert err:", err)
		trackUpload.Status = "FAILED"
		return
	}
	fmt.Println("Operation count:", opCount)
	trackUpload.DuplicateRecord = opCount.SkipCount
	if opCount.CreateCount == opCount.SkipCount {
		trackUpload.Status = "DUPLICATE"
	} else {
		trackUpload.Status = "COMPLETED"
	}

}

func insertFromFile(trackUpload *models.TrackUpload, forceCheck bool) (*OperationCount, error) {
	operationCount := &OperationCount{CreateCount: 0, SkipCount: 0, ForceCheck: forceCheck}
	app := GetApp()

	collection, _ := app.Dao().FindCollectionByNameOrId("track_upload")
	db, err := sql.Open("sqlite3", filepath.Join(app.DataDir(), "storage", collection.Id, trackUpload.BaseModel.Id, trackUpload.File))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var taskName string
	var beginDate, endDate types.DateTime

	err = db.QueryRow("select taskName, beginDate, endDate from TrackItems ORDER BY id ASC LIMIT 1;").Scan(&taskName, &beginDate, &endDate)

	if err != nil {
		return nil, err
	}

	err = db.QueryRow("select COUNT(*) FROM TrackItems;").Scan(&trackUpload.TotalRecord)

	if err != nil {
		return nil, err
	}

	if err := app.Dao().Save(trackUpload); err != nil {
		fmt.Println("Failed updating record")
		return nil, err
	}

	record, err := app.Dao().FindFirstRecordByFilter(
		"track_items", "user = {:user} && task_name = {:task_name} && source = {:source} && begin_date = {:begin_date} && end_date = {:end_date}",
		dbx.Params{"user": trackUpload.User,
			"task_name":  taskName,
			"source":     trackUpload.Source,
			"begin_date": beginDate,
			"end_date":   endDate})

	if err != nil {
		fmt.Println("No record found =", err)
	}
	var queryToExecute string
	queryCheckRequired := false
	if record == nil || operationCount.ForceCheck {
		queryToExecute = "select id, app, taskName, title, beginDate, endDate FROM TrackItems ORDER BY id ASC;"
		fmt.Println("No checks required")
	} else {
		queryToExecute = "select id, app, taskName, title, beginDate, endDate FROM TrackItems ORDER BY id DESC;"
		queryCheckRequired = true
		fmt.Println("Check required")
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {

		rows, err := db.Query(queryToExecute)
		if err != nil {
			fmt.Println("err =", err)
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {
			trackItems := &models.TrackItems{User: trackUpload.User, Source: trackUpload.Source}
			err = rows.Scan(&trackItems.TrackId, &trackItems.App, &trackItems.TaskName, &trackItems.Title, &trackItems.BeginDate, &trackItems.EndDate)

			if err != nil {
				fmt.Println("err =", err)
				log.Fatal(err)
			}

			if queryCheckRequired {
				record, err := app.Dao().FindFirstRecordByFilter(
					"track_items", "user = {:user} && task_name = {:task_name} && source = {:source} && begin_date = {:begin_date} && end_date = {:end_date}",
					dbx.Params{"user": trackUpload.User,
						"task_name":  trackItems.TaskName,
						"source":     trackItems.Source,
						"begin_date": trackItems.BeginDate,
						"end_date":   trackItems.EndDate})

				if err != nil {
					fmt.Println("Errorrrr =", err)
					return err
				}

				if record != nil {
					if operationCount.ForceCheck {
						operationCount.SkipCount += 1
						continue
					} else {
						break
					}
				}
			}

			if err := txDao.Save(trackItems); err != nil {
				return err
			}
			operationCount.CreateCount += 1
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return operationCount, nil
}

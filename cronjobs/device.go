package cronjobs

import (
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/shashank-sharma/backend/logger"
)

func TrackDevices(app *pocketbase.PocketBase) error {
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
		logger.Error.Println("No queries found")
	}
	for _, record := range records {
		record.Set("is_online", false)
		logger.Debug.Println("Working for record =", record)
		if err := app.Dao().SaveRecord(record); err != nil {
			logger.Error.Println("Error saving cron: ", err)
		}
	}
	return nil
}

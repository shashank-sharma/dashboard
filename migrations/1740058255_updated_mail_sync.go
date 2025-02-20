package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3850076299")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"hidden": false,
			"id": "select3548482499",
			"maxSelect": 1,
			"name": "sync_status",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"idle",
				"in_progress",
				"completed",
				"failed"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3850076299")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("select3548482499")

		return app.Save(collection)
	})
}

package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("650ed4q6e8zapgc")
		if err != nil {
			return err
		}

		// add
		new_last_online := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "95ckq7yd",
			"name": "last_online",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_last_online)
		collection.Schema.AddField(new_last_online)

		// add
		new_last_sync := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "ksebcvd0",
			"name": "last_sync",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_last_sync)
		collection.Schema.AddField(new_last_sync)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("650ed4q6e8zapgc")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("95ckq7yd")

		// remove
		collection.Schema.RemoveField("ksebcvd0")

		return dao.SaveCollection(collection)
	})
}

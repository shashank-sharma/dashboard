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
		new_sync_events := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "jcm4fjdv",
			"name": "sync_events",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_sync_events)
		collection.Schema.AddField(new_sync_events)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("650ed4q6e8zapgc")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("jcm4fjdv")

		return dao.SaveCollection(collection)
	})
}

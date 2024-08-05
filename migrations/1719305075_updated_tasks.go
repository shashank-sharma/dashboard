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

		collection, err := dao.FindCollectionByNameOrId("am2sdcobu2q3nfw")
		if err != nil {
			return err
		}

		// add
		new_project := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "fb5pbsb2",
			"name": "project",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_project)
		collection.Schema.AddField(new_project)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("am2sdcobu2q3nfw")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("fb5pbsb2")

		return dao.SaveCollection(collection)
	})
}

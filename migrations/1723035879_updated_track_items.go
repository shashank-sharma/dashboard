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

		collection, err := dao.FindCollectionByNameOrId("z60l06ij6ugtt49")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("l2r1wraq")

		// add
		new_device := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "fkrumvdd",
			"name": "device",
			"type": "relation",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "650ed4q6e8zapgc",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_device)
		collection.Schema.AddField(new_device)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z60l06ij6ugtt49")
		if err != nil {
			return err
		}

		// add
		del_source := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "l2r1wraq",
			"name": "source",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), del_source)
		collection.Schema.AddField(del_source)

		// remove
		collection.Schema.RemoveField("fkrumvdd")

		return dao.SaveCollection(collection)
	})
}

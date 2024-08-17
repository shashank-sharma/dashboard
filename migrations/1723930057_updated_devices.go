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
		new_is_public := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "orilnw0a",
			"name": "is_public",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_is_public)
		collection.Schema.AddField(new_is_public)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("650ed4q6e8zapgc")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("orilnw0a")

		return dao.SaveCollection(collection)
	})
}

package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "8h4vlopjsygl5yl",
			"created": "2024-12-22 12:58:38.153Z",
			"updated": "2024-12-22 12:58:38.153Z",
			"name": "track_focus",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "9vjnxfxj",
					"name": "user",
					"type": "relation",
					"required": true,
					"presentable": false,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": null
					}
				},
				{
					"system": false,
					"id": "ib9jlte0",
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
				},
				{
					"system": false,
					"id": "ifrdd2ez",
					"name": "tags",
					"type": "select",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"maxSelect": 9,
						"values": [
							"work",
							"entertainment",
							"gaming",
							"reading",
							"creative",
							"finance",
							"research",
							"writing",
							"learning"
						]
					}
				},
				{
					"system": false,
					"id": "bfcjmkdx",
					"name": "metadata",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "wskfbftl",
					"name": "begin_date",
					"type": "date",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": "",
						"max": ""
					}
				},
				{
					"system": false,
					"id": "ovxltuqt",
					"name": "end_date",
					"type": "date",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": "",
						"max": ""
					}
				}
			],
			"indexes": [],
			"listRule": "@request.auth.id = user",
			"viewRule": "@request.auth.id = user",
			"createRule": "@request.auth.id = user",
			"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
			"deleteRule": "@request.auth.id = user",
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("8h4vlopjsygl5yl")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}

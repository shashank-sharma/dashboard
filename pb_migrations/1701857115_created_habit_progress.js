migrate((db) => {
  const collection = new Collection({
    "id": "frl5rotrwv3eesj",
    "created": "2023-12-06 10:05:15.517Z",
    "updated": "2023-12-06 10:05:15.517Z",
    "name": "habit_progress",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "ru4ezxgg",
        "name": "habit",
        "type": "relation",
        "required": true,
        "unique": false,
        "options": {
          "collectionId": "t9sk3rywb3ap8vt",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": []
        }
      },
      {
        "system": false,
        "id": "rmto1bva",
        "name": "user",
        "type": "relation",
        "required": false,
        "unique": false,
        "options": {
          "collectionId": "_pb_users_auth_",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": []
        }
      },
      {
        "system": false,
        "id": "3gllrfeg",
        "name": "date",
        "type": "date",
        "required": true,
        "unique": false,
        "options": {
          "min": "",
          "max": ""
        }
      },
      {
        "system": false,
        "id": "ygy2luu5",
        "name": "status",
        "type": "select",
        "required": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "completed",
            "pending",
            "skipped"
          ]
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("frl5rotrwv3eesj");

  return dao.deleteCollection(collection);
})

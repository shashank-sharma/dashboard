migrate((db) => {
  const collection = new Collection({
    "id": "kcxn4npza3pe5w7",
    "created": "2023-12-08 08:09:59.869Z",
    "updated": "2023-12-08 08:09:59.869Z",
    "name": "food_log",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "y1zey4md",
        "name": "user",
        "type": "relation",
        "required": true,
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
        "id": "8z7aqaur",
        "name": "name",
        "type": "text",
        "required": true,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "nc3hahpp",
        "name": "field",
        "type": "file",
        "required": true,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "maxSize": 5242880,
          "mimeTypes": [],
          "thumbs": [],
          "protected": true
        }
      },
      {
        "system": false,
        "id": "e7vpmfz0",
        "name": "tag",
        "type": "select",
        "required": true,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "values": [
            "breakfast",
            "lunch",
            "dinner",
            "snack",
            "extra"
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
  const collection = dao.findCollectionByNameOrId("kcxn4npza3pe5w7");

  return dao.deleteCollection(collection);
})

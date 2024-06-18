migrate((db) => {
  const collection = new Collection({
    "id": "oduxdvd88nnn3i0",
    "created": "2023-12-06 11:28:34.292Z",
    "updated": "2023-12-06 11:28:34.292Z",
    "name": "daily_log",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "tdggxldt",
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
        "id": "g5mw43pj",
        "name": "summary",
        "type": "editor",
        "required": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "2mfocdla",
        "name": "score",
        "type": "number",
        "required": false,
        "unique": false,
        "options": {
          "min": 1,
          "max": 5
        }
      },
      {
        "system": false,
        "id": "fezgfzju",
        "name": "day",
        "type": "date",
        "required": false,
        "unique": false,
        "options": {
          "min": "",
          "max": ""
        }
      },
      {
        "system": false,
        "id": "cj8vulaj",
        "name": "bath",
        "type": "bool",
        "required": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "xd5g1m9p",
        "name": "feeling",
        "type": "text",
        "required": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
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
  const collection = dao.findCollectionByNameOrId("oduxdvd88nnn3i0");

  return dao.deleteCollection(collection);
})

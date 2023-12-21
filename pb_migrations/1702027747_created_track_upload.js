migrate((db) => {
  const collection = new Collection({
    "id": "9dpns8heuo4oukt",
    "created": "2023-12-08 09:29:07.002Z",
    "updated": "2023-12-08 09:29:07.002Z",
    "name": "track_upload",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "eo2l6bdp",
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
        "id": "ytzn6wgn",
        "name": "source",
        "type": "text",
        "required": false,
        "unique": false,
        "options": {
          "min": null,
          "max": null,
          "pattern": ""
        }
      },
      {
        "system": false,
        "id": "ezv92del",
        "name": "file",
        "type": "file",
        "required": false,
        "unique": false,
        "options": {
          "maxSelect": 1,
          "maxSize": 5242880,
          "mimeTypes": [],
          "thumbs": [],
          "protected": false
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
  const collection = dao.findCollectionByNameOrId("9dpns8heuo4oukt");

  return dao.deleteCollection(collection);
})

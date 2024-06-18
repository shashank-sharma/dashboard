migrate((db) => {
  const collection = new Collection({
    "id": "t9sk3rywb3ap8vt",
    "created": "2023-12-06 09:29:38.045Z",
    "updated": "2023-12-06 09:29:38.045Z",
    "name": "habits",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "avmxjods",
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
        "id": "cuzws3g1",
        "name": "type",
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
        "id": "edrekxyy",
        "name": "status",
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
        "id": "ccjvzri6",
        "name": "priority",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": null,
          "max": null
        }
      },
      {
        "system": false,
        "id": "qwb4x3gg",
        "name": "streak",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": null,
          "max": null
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
  const collection = dao.findCollectionByNameOrId("t9sk3rywb3ap8vt");

  return dao.deleteCollection(collection);
})

migrate((db) => {
  const collection = new Collection({
    "id": "w7w8875kihaao80",
    "created": "2023-12-06 10:54:54.166Z",
    "updated": "2023-12-06 10:54:54.166Z",
    "name": "life_balance",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "pgfnlwox",
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
        "id": "nkgw7pna",
        "name": "relationship",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "b4dtrsoz",
        "name": "health",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "rpotfkho",
        "name": "career",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "szl0uc5f",
        "name": "growth",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "n14uzvlu",
        "name": "life",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "udc4d4vc",
        "name": "social",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "opd7l8op",
        "name": "hobby",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
        }
      },
      {
        "system": false,
        "id": "erzuwcfq",
        "name": "finance",
        "type": "number",
        "required": true,
        "unique": false,
        "options": {
          "min": 1,
          "max": 10
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
  const collection = dao.findCollectionByNameOrId("w7w8875kihaao80");

  return dao.deleteCollection(collection);
})

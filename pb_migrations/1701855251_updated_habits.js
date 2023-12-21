migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("t9sk3rywb3ap8vt")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "djgqz8o0",
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
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("t9sk3rywb3ap8vt")

  // remove
  collection.schema.removeField("djgqz8o0")

  return dao.saveCollection(collection)
})

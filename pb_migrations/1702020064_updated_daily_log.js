migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oduxdvd88nnn3i0")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "fezgfzju",
    "name": "date",
    "type": "date",
    "required": true,
    "unique": false,
    "options": {
      "min": "",
      "max": ""
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("oduxdvd88nnn3i0")

  // update
  collection.schema.addField(new SchemaField({
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
  }))

  return dao.saveCollection(collection)
})

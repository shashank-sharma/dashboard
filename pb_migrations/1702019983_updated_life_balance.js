migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("w7w8875kihaao80")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "4xuu1iui",
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
  const collection = dao.findCollectionByNameOrId("w7w8875kihaao80")

  // remove
  collection.schema.removeField("4xuu1iui")

  return dao.saveCollection(collection)
})

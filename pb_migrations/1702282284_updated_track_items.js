migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("z60l06ij6ugtt49")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_0t9KXmL` ON `track_items` (\n  `begin_date`,\n  `end_date`,\n  `task_name`\n)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("z60l06ij6ugtt49")

  collection.indexes = [
    "CREATE INDEX `idx_0t9KXmL` ON `track_items` (\n  `begin_date`,\n  `end_date`,\n  `task_name`\n)"
  ]

  return dao.saveCollection(collection)
})

migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("frl5rotrwv3eesj")

  collection.listRule = "@request.auth.id = user"
  collection.viewRule = "@request.auth.id = user"
  collection.createRule = "@request.auth.id = user"
  collection.updateRule = "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)"
  collection.deleteRule = "@request.auth.id = user"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("frl5rotrwv3eesj")

  collection.listRule = null
  collection.viewRule = null
  collection.createRule = null
  collection.updateRule = null
  collection.deleteRule = null

  return dao.saveCollection(collection)
})
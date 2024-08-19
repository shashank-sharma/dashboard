package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("bu8wqeaacmnvshf")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.id = user")

		collection.ViewRule = types.Pointer("@request.auth.id = user")

		collection.CreateRule = types.Pointer("@request.auth.id = user")

		collection.UpdateRule = types.Pointer("@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)")

		collection.DeleteRule = types.Pointer("@request.auth.id = user")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("bu8wqeaacmnvshf")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = nil

		collection.CreateRule = nil

		collection.UpdateRule = nil

		collection.DeleteRule = nil

		return dao.SaveCollection(collection)
	})
}

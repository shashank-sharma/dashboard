package store

import (
	"sync"

	"github.com/pocketbase/pocketbase/daos"
)

var (
	once sync.Once
	dao  *daos.Dao
)

func InitDao(newDao *daos.Dao) {
	once.Do(func() {
		dao = newDao
	})
}

func GetDao() *daos.Dao {
	if dao == nil {
		panic("dao has not been initialized")
	}
	return dao
}

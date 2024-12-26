package store

import (
	"sync"

	"github.com/pocketbase/pocketbase"
)

var (
	once sync.Once
	dao  *pocketbase.PocketBase
)

func InitApp(newApp *pocketbase.PocketBase) {
	once.Do(func() {
		dao = newApp
	})
}

func GetDao() *pocketbase.PocketBase {
	if dao == nil {
		panic("dao has not been initialized")
	}
	return dao
}

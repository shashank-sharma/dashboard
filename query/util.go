package query

import (
	"github.com/pocketbase/dbx"
)

func structToHashExp(filterStruct map[string]interface{}) dbx.HashExp {
	hashExp := dbx.HashExp{}
	for key, value := range filterStruct {
		hashExp[key] = value
	}
	return hashExp
}

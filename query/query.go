package query

import (
	"github.com/pocketbase/dbx"
	"github.com/shashank-sharma/backend/logger"
	"github.com/shashank-sharma/backend/models"
	"github.com/shashank-sharma/backend/store"
)

func BaseModelQuery[T models.Model](m T) *dbx.SelectQuery {
	return store.GetDao().ModelQuery(m)
}

func FindById[T models.Model](id string) (T, error) {
	var m T
	query := BaseModelQuery(m)
	err := query.
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(&m)

	if err != nil {
		logger.Debug.Println(err)
		var zeroValue T
		return zeroValue, err
	}

	return m, nil
}

func FindByFilter[T models.Model](filterStruct map[string]interface{}) (T, error) {
	var m T
	query := BaseModelQuery(m)
	filter := structToHashExp(filterStruct)

	err := query.
		AndWhere(filter).
		Limit(1).
		One(&m)

	if err != nil {
		var zeroValue T
		return zeroValue, err
	}

	return m, nil
}

func SaveRecord(model models.Model) error {
	if err := store.GetDao().Save(model); err != nil {
		return err
	}
	return nil
}

func UpsertRecord[T models.Model](model T, filterStruct map[string]interface{}) error {
	record, err := FindByFilter[T](filterStruct)
	if err == nil {
		model.SetId(record.GetId())
	}
	if err := SaveRecord(model); err != nil {
		return err
	}
	return nil
}

func UpdateRecord[T models.Model](filterId string, updateStruct map[string]interface{}) error {
	var m T
	record, err := store.GetDao().FindRecordById(m.TableName(), filterId)
	if err != nil {
		return err
	}

	for key, value := range updateStruct {
		record.Set(key, value)
	}

	if err := store.GetDao().SaveRecord(record); err != nil {
		return err
	}
	return nil
}

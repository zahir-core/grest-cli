package app

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"

	"gorm.io/gorm"
	"grest.dev/grest"
)

type dbQuery struct {
	grest.DBQuery
}

func First(db *gorm.DB, model ModelInterface, query url.Values) error {
	query.Add(grest.QueryInclude, "all")
	res, err := Find(db, model, query)
	if err != nil {
		return NewError(http.StatusInternalServerError, err.Error())
	}
	if len(res) == 0 {
		return gorm.ErrRecordNotFound
	}
	b, err := json.Marshal(res[0])
	if err != nil {
		return NewError(http.StatusInternalServerError, err.Error())
	}
	err = json.Unmarshal(b, model)
	if err != nil {
		return NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func Find(db *gorm.DB, model ModelInterface, query url.Values) ([]map[string]any, error) {
	q := &dbQuery{}
	q.DB = db
	q.Schema = model.GetSchema()
	q.Query = query
	return q.Find(q.Schema, query)
}

func PaginationInfo(db *gorm.DB, model ModelInterface, query url.Values) (int64, int, int, int, error) {
	var err error
	count, page, perPage, pageCount := int64(0), 0, 0, 0
	if query.Get(grest.QueryDisablePagination) == "true" {
		return count, page, perPage, pageCount, err
	}

	q := &dbQuery{}
	q.DB = db
	q.Schema = model.GetSchema()
	q.Query = query
	tx, err := q.Prepare(db, q.Schema, query)
	if err != nil {
		return count, page, perPage, pageCount, err
	}
	err = tx.Count(&count).Error
	if err != nil || query.Get(grest.QueryLimit) == "0" {
		return count, page, perPage, pageCount, err
	}
	page, perPage = q.GetPageLimit()
	pageCount = int(math.Ceil(float64(count) / float64(perPage)))
	return count, page, perPage, pageCount, err
}

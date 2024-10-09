package codegentemplate

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"grest.dev/cmd/codegentemplate/app"
)

// UseCase returns a useCase for expected use case functional.
func UseCase(ctx app.Ctx, query ...url.Values) useCase {
	u := useCase{
		Ctx:   &ctx,
		Query: url.Values{},
	}
	if len(query) > 0 {
		u.Query = query[0]
	}
	return u
}

// useCase provides a convenient interface for codegentemplate use case, use UseCase to access useCase.
type useCase struct {

	// injectable dependencies
	Ctx   *app.Ctx   `json:"-" db:"-" gorm:"-"`
	Query url.Values `json:"-" db:"-" gorm:"-"`
}

// Async return useCase for async process
func (u useCase) Async(ctx app.Ctx, query ...url.Values) useCase {
	ctx.IsAsync = true
	return UseCase(ctx, query...)
}

// GetByID returns the codegentemplate data for the specified ID.
func (u useCase) GetByID(id string) (CodeGenTemplate, error) {
	res := CodeGenTemplate{}

	// check permission
	err := u.Ctx.ValidatePermission("end_point.detail")
	if err != nil {
		return res, err
	}

	// validate if codegentemplate is exists
	key := "id"
	if !app.Validator().IsValid(id, "uuid") {
		key = "code"
	}
	realID, err := u.GetIDByKey(key, id)

	// get from cache and return if exists
	cacheKey := CodeGenTemplate{}.EndPoint() + "." + id
	app.Cache().Get(cacheKey, &res)
	if res.ID.Valid {
		return res, err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return res, app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// get from db
	u.Query.Add("id", realID.String)
	err = app.Query().First(tx, &res, u.Query)
	if err != nil {
		return res, u.Ctx.NotFoundError(err, CodeGenTemplate{}.EndPoint(), key, id)
	}

	// save to cache and return if exists
	app.Cache().Set(cacheKey, res)
	return res, err
}

// Get returns the list of codegentemplate data.
func (u useCase) Get() (app.ListModel, error) {
	res := app.ListModel{}

	// check permission
	err := u.Ctx.ValidatePermission("end_point.list")
	if err != nil {
		return res, err
	}
	// get from cache and return if exists
	cacheKey := CodeGenTemplate{}.EndPoint() + "?" + u.Query.Encode()
	err = app.Cache().Get(cacheKey, &res)
	if err == nil {
		return res, err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return res, app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// set pagination info
	res.Count,
		res.PageContext.Page,
		res.PageContext.PerPage,
		res.PageContext.PageCount,
		err = app.Query().PaginationInfo(tx, &CodeGenTemplate{}, u.Query)
	if err != nil {
		return res, app.Error().New(http.StatusInternalServerError, err.Error())
	}
	// return data count if $per_page set to 0
	if res.PageContext.PerPage == 0 {
		return res, err
	}

	// find data
	data, err := app.Query().Find(tx, &CodeGenTemplate{}, u.Query)
	if err != nil {
		return res, app.Error().New(http.StatusInternalServerError, err.Error())
	}
	res.SetData(data, u.Query)

	// save to cache and return if exists
	app.Cache().Set(cacheKey, res)
	return res, err
}

// Create creates a new data codegentemplate with specified parameters.
func (u useCase) Create(param *CodeGenTemplate, paramCreate *ParamCreate) error {

	// check permission
	err := u.Ctx.ValidatePermission("end_point.create")
	if err != nil {
		return err
	}

	// validate param
	err = u.Ctx.ValidateParam(paramCreate)
	if err != nil {
		return err
	}

	// set default value for undefined field
	old := CodeGenTemplate{}
	err = u.setDefaultValue(old, param)
	if err != nil {
		return err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// save data to db
	err = tx.Model(param).Create(&param).Error
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// invalidate cache
	app.Cache().Invalidate(CodeGenTemplate{}.EndPoint())

	// save history (user activity), send webhook, etc
	go u.Ctx.Hook("POST", "create", param.ID.String, param)
	return nil
}

// UpdateByID updates the codegentemplate data for the specified ID with specified parameters.
func (u useCase) UpdateByID(id string, param *CodeGenTemplate, paramUpdate *ParamUpdate) error {

	// check permission
	err := u.Ctx.ValidatePermission("end_point.edit")
	if err != nil {
		return err
	}

	// validate param
	err = u.Ctx.ValidateParam(paramUpdate)
	if err != nil {
		return err
	}

	// get previous data
	old, err := u.GetByID(id)
	if err != nil {
		return err
	}

	// set default value for undefined field
	if err := u.setDefaultValue(old, param); err != nil {
		return err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// update data on the db
	err = tx.Model(param).Where("id = ?", old.ID).Updates(param).Error
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// invalidate cache
	app.Cache().Invalidate(CodeGenTemplate{}.EndPoint(), old.ID.String)

	// save history (user activity), send webhook, etc
	go u.Ctx.Hook("PUT", paramUpdate.Reason.String, old.ID.String, old)
	return nil
}

// PartiallyUpdateByID updates the codegentemplate data for the specified ID with specified parameters.
func (u useCase) PartiallyUpdateByID(id string, param *CodeGenTemplate, paramUpdate *ParamPartiallyUpdate) error {

	// check permission
	err := u.Ctx.ValidatePermission("end_point.edit")
	if err != nil {
		return err
	}

	// validate param
	err = u.Ctx.ValidateParam(paramUpdate)
	if err != nil {
		return err
	}

	// get previous data
	old, err := u.GetByID(id)
	if err != nil {
		return err
	}

	// set default value for undefined field
	if err := u.setDefaultValue(old, param); err != nil {
		return err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// update data on the db
	err = tx.Model(param).Where("id = ?", old.ID).Updates(param).Error
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// invalidate cache
	app.Cache().Invalidate(CodeGenTemplate{}.EndPoint(), old.ID.String)

	// save history (user activity), send webhook, etc
	go u.Ctx.Hook("PATCH", paramUpdate.Reason.String, old.ID.String, old)
	return nil
}

// DeleteByID deletes the codegentemplate data for the specified ID.
func (u useCase) DeleteByID(id string, paramDelete *ParamDelete) error {

	// check permission
	err := u.Ctx.ValidatePermission("end_point.delete")
	if err != nil {
		return err
	}

	// validate param
	err = u.Ctx.ValidateParam(paramDelete)
	if err != nil {
		return err
	}

	// get previous data
	old, err := u.GetByID(id)
	if err != nil {
		return err
	}

	// prepare db for current ctx
	tx, err := u.Ctx.DB()
	if err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// update data on the db
	if err = tx.Model(paramDelete).Where("id = ?", old.ID).Update("deleted_at", time.Now().UTC()).Error; err != nil {
		return app.Error().New(http.StatusInternalServerError, err.Error())
	}

	// invalidate cache
	app.Cache().Invalidate(CodeGenTemplate{}.EndPoint(), old.ID.String)

	// save history (user activity), send webhook, etc
	go u.Ctx.Hook("DELETE", paramDelete.Reason.String, old.ID.String, old)
	return nil
}

// GetIDByKey get codegentemplate id by unique key.
func (u useCase) GetIDByKey(key, val string) (app.NullUUID, error) {
	d := &CodeGenTemplate{}
	app.Cache().Get(CodeGenTemplate{}.EndPoint()+"."+val, d)
	if d.ID.String != "" {
		return d.ID, nil
	}
	tx, err := u.Ctx.DB()
	if err != nil {
		return d.ID, app.Error().New(http.StatusInternalServerError, err.Error())
	}
	fKey := "id"
	if key != "id" {
		fKey = "code"
	}
	tx.Model(d).Where(fKey+" = ?", val).Where("deleted_at is null").Take(d)
	if d.ID.String != "" {
		return d.ID, nil
	}
	return d.ID, u.Ctx.NotFoundError(nil, CodeGenTemplate{}.EndPoint(), key, val)
}

// ParseParamCreate return *ParamCreate from CodeGenTemplate
func (u useCase) ParseParamCreate(data CodeGenTemplate) *ParamCreate {
	param := &ParamCreate{}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, param)
	return param
}

// setDefaultValue set default value of undefined field when create or update codegentemplate data.
func (u useCase) setDefaultValue(old CodeGenTemplate, new *CodeGenTemplate) error {
	if !old.ID.Valid {
		new.ID = app.NewNullUUID()
	} else {
		new.ID = old.ID
	}
	return nil
}

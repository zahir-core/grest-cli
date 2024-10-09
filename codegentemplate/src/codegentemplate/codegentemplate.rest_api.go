package codegentemplate

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"grest.dev/cmd/codegentemplate/app"
)

// REST returns a *restAPI.
func REST() *restAPI {
	return &restAPI{}
}

// restAPI provides a convenient interface for codegentemplate REST API handler.
type restAPI struct {
	UseCase useCase
}

// injectDeps inject the dependencies of the codegentemplate REST API handler.
func (r *restAPI) injectDeps(c *fiber.Ctx) error {
	ctx, ok := c.Locals(app.CtxKey).(*app.Ctx)
	if !ok {
		return app.Error().New(http.StatusInternalServerError, "ctx is not found")
	}
	r.UseCase = UseCase(*ctx, app.Query().Parse(c.OriginalURL()))
	return nil
}

// GetByID is the REST API handler for `GET /api/v3/end_point/{id}`.
func (r *restAPI) GetByID(c *fiber.Ctx) error {
	if err := r.injectDeps(c); err != nil {
		return app.Server().Error(c, err)
	}
	res, err := r.UseCase.GetByID(c.Params("id"))
	if err != nil {
		return app.Server().Error(c, err)
	}
	return c.JSON(app.Query().Return(res, res.IsFlat()))
}

// Get is the REST API handler for `GET /api/v3/end_point`.
func (r *restAPI) Get(c *fiber.Ctx) error {
	if err := r.injectDeps(c); err != nil {
		return app.Server().Error(c, err)
	}
	res, err := r.UseCase.Get()
	if err != nil {
		return app.Server().Error(c, err)
	}
	res.SetLink(c)
	model := &CodeGenTemplate{}
	return c.JSON(app.Query().Return(res, model.IsFlat()))
}

// Create is the REST API handler for `POST /api/v3/end_point`.
func (r *restAPI) Create(c *fiber.Ctx) error {
	if err := r.injectDeps(c); err != nil {
		return app.Server().Error(c, err)
	}
	param := &CodeGenTemplate{}
	paramCreate := &ParamCreate{}
	if err := app.Query().BindJSON(c.Body(), param, paramCreate); err != nil {
		return app.Server().Error(c, app.Error().New(http.StatusBadRequest, err.Error()))
	}
	if err := r.UseCase.Create(param, paramCreate); err != nil {
		return app.Server().Error(c, err)
	}
	if r.UseCase.Query.Get("is_skip_return") == "true" {
		return c.Status(http.StatusCreated).JSON(map[string]any{"message": "Success"})
	}
	res, err := r.UseCase.GetByID(param.ID.String)
	if err != nil {
		return app.Server().Error(c, err)
	}
	return c.Status(http.StatusCreated).JSON(app.Query().Return(res, res.IsFlat()))
}

// UpdateByID is the REST API handler for `PUT /api/v3/end_point/{id}`.
func (r *restAPI) UpdateByID(c *fiber.Ctx) error {
	if err := r.injectDeps(c); err != nil {
		return app.Server().Error(c, err)
	}
	param := &CodeGenTemplate{}
	paramUpdate := &ParamUpdate{}
	if err := app.Query().BindJSON(c.Body(), param, paramUpdate); err != nil {
		return app.Server().Error(c, app.Error().New(http.StatusBadRequest, err.Error()))
	}
	if err := r.UseCase.UpdateByID(c.Params("id"), param, paramUpdate); err != nil {
		return app.Server().Error(c, err)
	}
	if r.UseCase.Query.Get("is_skip_return") == "true" {
		return c.JSON(map[string]any{"message": "Success"})
	}
	res, err := r.UseCase.GetByID(c.Params("id"))
	if err != nil {
		return app.Server().Error(c, err)
	}
	return c.JSON(app.Query().Return(res, res.IsFlat()))
}

// PartiallyUpdateByID is the REST API handler for `PATCH /api/v3/end_point/{id}`.
func (r *restAPI) PartiallyUpdateByID(c *fiber.Ctx) error {
	if err := r.injectDeps(c); err != nil {
		return app.Server().Error(c, err)
	}
	param := &CodeGenTemplate{}
	paramUpdate := &ParamPartiallyUpdate{}
	if err := app.Query().BindJSON(c.Body(), param, paramUpdate); err != nil {
		return app.Server().Error(c, app.Error().New(http.StatusBadRequest, err.Error()))
	}
	if err := r.UseCase.PartiallyUpdateByID(c.Params("id"), param, paramUpdate); err != nil {
		return app.Server().Error(c, err)
	}
	if r.UseCase.Query.Get("is_skip_return") == "true" {
		return c.JSON(map[string]any{"message": "Success"})
	}
	res, err := r.UseCase.GetByID(c.Params("id"))
	if err != nil {
		return app.Server().Error(c, err)
	}
	return c.JSON(app.Query().Return(res, res.IsFlat()))
}

// DeleteByID is the REST API handler for `DELETE /api/v3/end_point/{id}`.
func (r *restAPI) DeleteByID(c *fiber.Ctx) error {
	err := r.injectDeps(c)
	if err != nil {
		return app.Server().Error(c, err)
	}
	id := c.Params("id")
	paramDelete := &ParamDelete{}
	if err := app.Query().BindJSON(c.Body(), paramDelete); err != nil {
		return app.Server().Error(c, app.Error().New(http.StatusBadRequest, err.Error()))
	}
	if err = r.UseCase.DeleteByID(id, paramDelete); err != nil {
		return app.Server().Error(c, err)
	}
	return c.JSON(r.UseCase.Ctx.Deleted(CodeGenTemplate{}.EndPoint(), "id", id))
}

package codegentemplate

import "grest.dev/cmd/codegentemplate/app"

// CodeGenTemplate is the main model of CodeGenTemplate data. It provides a convenient interface for app.ModelInterface
type CodeGenTemplate struct {
	app.Model
	ID               app.NullUUID     `json:"id"                 db:"m.id"                 gorm:"column:id;primaryKey"`
	CreatedUserID    app.NullUUID     `json:"created.user.id"    db:"m.created_by_user_id" gorm:"column:created_by_user_id;index:table_name_created_by_user_id"`
	CreatedUserName  app.NullString   `json:"created.user.name"  db:"cu.name"              gorm:"-"`
	CreatedUserEmail app.NullString   `json:"created.user.email" db:"cu.nama"              gorm:"-"`
	CreatedAt        app.NullDateTime `json:"created.time"       db:"m.created_at"         gorm:"column:created_at"`
	UpdatedUserID    app.NullUUID     `json:"updated.user.id"    db:"m.updated_by_user_id" gorm:"column:updated_by_user_id;index:table_name_updated_by_user_id"`
	UpdatedUserName  app.NullString   `json:"updated.user.name"  db:"uu.name"              gorm:"-"`
	UpdatedUserEmail app.NullString   `json:"updated.user.email" db:"uu.nama"              gorm:"-"`
	UpdatedAt        app.NullDateTime `json:"updated.time"       db:"m.updated_at"         gorm:"column:updated_at"`
	DeletedAt        app.NullDateTime `json:"deleted.time"       db:"m.deleted_at"         gorm:"column:deleted_at"`
}

// EndPoint returns the CodeGenTemplate end point, it used for cache key, etc.
func (CodeGenTemplate) EndPoint() string {
	return "end_point"
}

// TableVersion returns the versions of the CodeGenTemplate table in the database.
// Change this value with date format YY.MM.DDHHii when any table structure changes.
func (CodeGenTemplate) TableVersion() string {
	return "28.06.291152"
}

// TableName returns the name of the CodeGenTemplate table in the database.
func (CodeGenTemplate) TableName() string {
	return "table_name"
}

// TableAliasName returns the table alias name of the CodeGenTemplate table, used for querying.
func (CodeGenTemplate) TableAliasName() string {
	return "m"
}

// GetRelations returns the relations of the CodeGenTemplate data in the database, used for querying.
func (m *CodeGenTemplate) GetRelations() map[string]map[string]any {
	m.AddRelation("left", "sistem", "cu", []map[string]any{{"column1": "cu.id", "column2": "m.created_by_user_id"}})
	m.AddRelation("left", "sistem", "uu", []map[string]any{{"column1": "uu.id", "column2": "m.updated_by_user_id"}})
	return m.Relations
}

// GetFilters returns the filter of the CodeGenTemplate data in the database, used for querying.
func (m *CodeGenTemplate) GetFilters() []map[string]any {
	m.AddFilter(map[string]any{"column1": "m.deleted_at", "operator": "=", "value": nil})
	return m.Filters
}

// GetSorts returns the default sort of the CodeGenTemplate data in the database, used for querying.
func (m *CodeGenTemplate) GetSorts() []map[string]any {
	m.AddSort(map[string]any{"column": "m.updated_at", "direction": "desc"})
	return m.Sorts
}

// GetFields returns list of the field of the CodeGenTemplate data in the database, used for querying.
func (m *CodeGenTemplate) GetFields() map[string]map[string]any {
	m.SetFields(m)
	return m.Fields
}

// GetSchema returns the CodeGenTemplate schema, used for querying.
func (m *CodeGenTemplate) GetSchema() map[string]any {
	return m.SetSchema(m)
}

// OpenAPISchemaName returns the name of the CodeGenTemplate schema in the open api documentation.
func (CodeGenTemplate) OpenAPISchemaName() string {
	return "SchemaCategory.CodeGenTemplate"
}

// ParamCreate is the expected parameters for create a new CodeGenTemplate data.
type ParamCreate struct {
	UseCaseHandler
}

// ParamUpdate is the expected parameters for update the CodeGenTemplate data.
type ParamUpdate struct {
	UseCaseHandler
	Reason app.NullString `json:"reason" gorm:"-" validate:"required"`
}

// ParamPartiallyUpdate is the expected parameters for partially update the CodeGenTemplate data.
type ParamPartiallyUpdate struct {
	UseCaseHandler
	Reason app.NullString `json:"reason" gorm:"-" validate:"required"`
}

// ParamDelete is the expected parameters for delete the CodeGenTemplate data.
type ParamDelete struct {
	UseCaseHandler
	Reason app.NullString `json:"reason" gorm:"-" validate:"required"`
}

package tests

import "grest.dev/grest"

type Book struct {
	ID          grest.NullUUID     `json:"id"           db:"c.id"`
	Code        grest.NullString   `json:"code"         db:"c.code"`
	Name        grest.NullString   `json:"name"         db:"c.name"`
	AuthorID    grest.NullUUID     `json:"author.id"    db:"c.author_id"`
	AuthorName  grest.NullString   `json:"author.name"  db:"u.name"`
	AuthorEmail grest.NullString   `json:"author.email" db:"u.email"`
	CreatedAt   grest.NullDateTime `json:"created.time" db:"c.created_at"`
	UpdatedAt   grest.NullDateTime `json:"updated.time" db:"c.updated_at"`
	DeletedAt   grest.NullDateTime `json:"-"            db:"c.deleted_at"`
	Foo         grest.String
	Bar         grest.String
	Baz         grest.String
}

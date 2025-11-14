package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Image holds the schema definition for the Images entity.
type Image struct {
	ent.Schema
}

// Fields of the Images.
func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("imageUrl"),
	}
}

// Edges of the Images.
func (Image) Edges() []ent.Edge {
	return nil
}

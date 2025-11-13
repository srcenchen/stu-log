package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Score holds the schema definition for the Score entity.
type Score struct {
	ent.Schema
}

// Fields of the Score.
func (Score) Fields() []ent.Field {
	return []ent.Field{
		field.Int("score"),
	}
}

// Edges of the Score.
func (Score) Edges() []ent.Edge {
	return nil
}

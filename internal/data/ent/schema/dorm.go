package schema

import "entgo.io/ent"

// Dorm holds the schema definition for the Dorm entity.
type Dorm struct {
	ent.Schema
}

// Fields of the Dorm.
func (Dorm) Fields() []ent.Field {
	return nil
}

// Edges of the Dorm.
func (Dorm) Edges() []ent.Edge {
	return nil
}

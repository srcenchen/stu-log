package schema

import "entgo.io/ent"

// Logs holds the schema definition for the Logs entity.
type Logs struct {
	ent.Schema
}

// Fields of the Logs.
func (Logs) Fields() []ent.Field {
	return nil
}

// Edges of the Logs.
func (Logs) Edges() []ent.Edge {
	return nil
}

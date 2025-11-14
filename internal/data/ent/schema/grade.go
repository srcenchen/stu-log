package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Grade holds the schema definition for the Grade entity.
type Grade struct {
	ent.Schema
}

// Fields of the Grade.
func (Grade) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("gradeName"),
	}
}

// Edges of the Grade.
func (Grade) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stuLogs", StuLog.Type).Ref("grade"),
		edge.From("class", Class.Type).Ref("grade"),
		edge.From("student", Student.Type).Ref("grade"),
		edge.From("dorm", Dorm.Type).Ref("grade"),
	}
}

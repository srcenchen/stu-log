package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Dorm holds the schema definition for the Dorm entity.
type Dorm struct {
	ent.Schema
}

// Fields of the Dorm.
func (Dorm) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("building"),
		field.String("dormNum"),
		field.String("sex"),
	}
}

// Edges of the Dorm.
func (Dorm) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("student", Student.Type).Ref("dorm"),
		edge.To("grade", Grade.Type).Unique().Required(),
		edge.From("stuLogs", StuLog.Type).Ref("dorm"),
	}
}

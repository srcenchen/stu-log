package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Class holds the schema definition for the Class entity.
type Class struct {
	ent.Schema
}

// Fields of the Class.
func (Class) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("className"),
	}
}

func (Class) Edges() []ent.Edge {
	return []ent.Edge{
		// 每个 Class 属于一个 Grade
		edge.To("grade", Grade.Type).Unique().Required(),
		edge.From("student", Student.Type).Ref("class"),
	}
}

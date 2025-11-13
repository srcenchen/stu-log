package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Student holds the schema definition for the Student entity.
type Student struct {
	ent.Schema
}

// Fields of the Student.
func (Student) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("学生姓名"),
		field.String("sex").NotEmpty().Comment("性别"),
		field.Int("score").Default(100).StorageKey("student_score"),
	}
}

// Edges of the Student.
func (Student) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("grade", Grade.Type).Unique().Required(),
		edge.To("class", Class.Type).Unique().Required(),
		edge.To("dorm", Dorm.Type).Unique(),
	}
}

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
		field.Int64("id"),
		field.String("name").NotEmpty().Comment("学生姓名"),
		field.String("stuNum").NotEmpty().Unique().Comment("学号"),
		field.String("sex").NotEmpty().Comment("性别"),
		field.Int32("score").Default(100).Comment("学生分数"),
		field.String("dormPos").Default("").Comment("宿舍床铺的位置"),
	}
}

// Edges of the Student.
func (Student) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("grade", Grade.Type).Unique().Required(),
		edge.To("class", Class.Type).Unique().Required(),
		edge.To("dorm", Dorm.Type).Unique(),
		edge.To("stuLogs", StuLog.Type),
	}
}

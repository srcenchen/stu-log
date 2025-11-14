package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StuLog holds the schema definition for the Logs entity.
type StuLog struct {
	ent.Schema
}

// Fields of the Logs.
func (StuLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("content").Comment("违纪备注"),
		field.Bool("revoked").Default(false).Comment("是否被撤销"),
		field.Int32("score").Comment("加减分情况"),
		field.Time("time").Comment("违纪时间"),
	}
}

// Edges of the Logs.
func (StuLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("class", Class.Type),
		edge.To("grade", Grade.Type).Required(),
		edge.To("rule", Rule.Type).Required(),
		edge.To("students", Student.Type),
		edge.To("images", Image.Type),
	}
}

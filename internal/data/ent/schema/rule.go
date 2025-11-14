package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Rule holds the schema definition for the Rules entity.
type Rule struct {
	ent.Schema
}

// Fields of the Rules.
func (Rule) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("content").Comment("内容").Unique(),
		field.Int32("score").Comment("分数，加分还是扣分"),
		field.String("group").Comment("分组"),
	}
}

// Edges of the Rules.
func (Rule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stuLogs", StuLog.Type).Ref("rule"),
	}
}

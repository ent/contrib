package entproto

import (
	"entgo.io/ent"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"github.com/mitchellh/mapstructure"
)

type FieldGroup struct {
	Name   string
	Fields []ent.Field
}

// type FieldGroups []FieldGroup

// func (fg FieldGroups) Fields() []ent.Field {
// 	fields := []ent.Field{}
// 	for _, g := range fg {
// 		fields = append(fields, g.Fields...)
// 	}
// 	return fields
// }

func (fg FieldGroups) ByName(names ...string) FieldGroups {
	result := FieldGroups{
		groups: []*FieldGroup{},
		Names:  []string{},
	}
	for _, g := range fg.groups {
		for _, name := range names {
			if g.Name == name {
				result.groups = append(result.groups, g)
				result.Names = append(result.Names, g.Name)
			}
		}
	}
	return result
}

type FieldGroups struct {
	groups []*FieldGroup
	Names  []string
}

func Groups() *FieldGroups {
	return &FieldGroups{
		groups: []*FieldGroup{},
		Names:  []string{},
	}
}

func (g *FieldGroups) Group(name string, builder func(*FieldGroup)) *FieldGroups {
	f := &FieldGroup{
		Name: name,
	}
	builder(f)

	// Append the annotations
	for _, f := range f.Fields {
		d := f.Descriptor()
		d.Annotations = append(d.Annotations, GroupName(name))
	}

	g.groups = append(g.groups, f)
	g.Names = append(g.Names, name)

	return g
}

func (fg FieldGroups) Fields() []ent.Field {
	fields := []ent.Field{}
	for _, g := range fg.groups {
		fields = append(fields, g.Fields...)
	}
	return fields
}

const GroupNameAnnotation = "ProtoFieldGroupName"

func GroupName(name string) schema.Annotation {
	f := &pbgroupName{GroupName: name}
	return f
}

type pbgroupName struct {
	GroupName string
}

func (f pbgroupName) Name() string {
	return GroupNameAnnotation
}

func extractGroupNameAnnotation(sch *gen.Field) *pbgroupName {
	annot, ok := sch.Annotations[GroupNameAnnotation]
	if !ok {
		return nil
	}

	var groupName pbgroupName
	mapstructure.Decode(annot, &groupName)

	return &groupName
}

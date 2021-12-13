package create

import (
	"github.com/PM-Master/policy-machine-go/policy"
)

type create struct {
	name       string
	kind       policy.Kind
	properties map[string]string
}

func (c *create) In(parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       c.name,
		Kind:       c.kind,
		Properties: c.properties,
		Parents:    parents,
	}
}

func (c *create) WithProperties(props ...string) *create {
	properties := make(map[string]string)
	for i := 0; i < len(props); i += 2 {
		properties[props[i]] = props[i+1]
	}

	c.properties = properties

	return c
}

func UserAttribute(name string) *create {
	return &create{
		name:       name,
		kind:       policy.UserAttribute,
		properties: make(map[string]string),
	}
}

func User(name string) *create {
	return &create{
		name:       name,
		kind:       policy.User,
		properties: make(map[string]string),
	}
}

func ObjectAttribute(name string) *create {
	return &create{
		name:       name,
		kind:       policy.ObjectAttribute,
		properties: make(map[string]string),
	}
}

func Object(name string) *create {
	return &create{
		name:       name,
		kind:       policy.Object,
		properties: make(map[string]string),
	}
}

func PolicyClass(name string) *policy.CreatePolicyStatement {
	return &policy.CreatePolicyStatement{
		Name: name,
	}
}

package deassign

import "github.com/PM-Master/policy-machine-go/policy"

type deassign struct {
	child   string
	parents []string
}

func UserAttribute(name string) *deassign {
	return &deassign{child: name}
}

func User(name string) *deassign {
	return &deassign{child: name}
}

func ObjectAttribute(name string) *deassign {
	return &deassign{child: name}
}

func Object(name string) *deassign {
	return &deassign{child: name}
}

func (a *deassign) From(parent string, parents ...string) *policy.DeassignStatement {
	parents = append(parents, parent)
	return &policy.DeassignStatement{
		Child:   a.child,
		Parents: parents,
	}
}

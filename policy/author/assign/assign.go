package assign

import "github.com/PM-Master/policy-machine-go/policy"

type assign struct {
	child   string
	parents []string
}

func UserAttribute(name string) *assign {
	return &assign{child: name}
}

func User(name string) *assign {
	return &assign{child: name}
}

func ObjectAttribute(name string) *assign {
	return &assign{child: name}
}

func Object(name string) *assign {
	return &assign{child: name}
}

func (a *assign) To(parent string, parents ...string) policy.AssignStatement {
	parents = append(parents, parent)
	return policy.AssignStatement{
		Child:   a.child,
		Parents: parents,
	}
}

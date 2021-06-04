package graph

import "fmt"

type (
	Node struct {
		Name       string            `json:"name"`
		Kind       Kind              `json:"kind"`
		Properties map[string]string `json:"properties"`
	}

	Operations map[string]bool

	Kind int
)

const (
	PolicyClass Kind = iota
	ObjectAttribute
	UserAttribute
	Object
	User

	AllOps = "*"
)

func (k Kind) String() string {
	switch k {
	case PolicyClass:
		return "PC"
	case ObjectAttribute:
		return "OA"
	case UserAttribute:
		return "UA"
	case Object:
		return "O"
	case User:
		return "U"
	default:
		return "nil"
	}
}

var (
	validAssignments = map[Kind]map[Kind]bool{
		PolicyClass:     {},
		ObjectAttribute: {ObjectAttribute: true, PolicyClass: true},
		UserAttribute:   {UserAttribute: true, PolicyClass: true},
		Object:          {ObjectAttribute: true},
		User:            {UserAttribute: true},
	}

	validAssociations = map[Kind]map[Kind]bool{
		PolicyClass:     {},
		ObjectAttribute: {},
		UserAttribute:   {ObjectAttribute: true, UserAttribute: true},
		Object:          {},
		User:            {},
	}
)

func CheckAssignment(childKind Kind, parentKind Kind) error {
	if !validAssignments[childKind][parentKind] {
		return fmt.Errorf("invalid assignment: %q to %q", childKind.String(), parentKind.String())
	}

	return nil
}

func CheckAssociation(subjectKind Kind, targetKind Kind) error {
	if !validAssociations[subjectKind][targetKind] {
		return fmt.Errorf("invalid association: %q to %q", subjectKind.String(), targetKind.String())
	}

	return nil
}

func ToOps(ops ...string) (operations Operations) {
	operations = make(map[string]bool)
	for _, op := range ops {
		operations[op] = true
	}
	return operations
}

func (o Operations) Contains(op string) bool {
	return o[op] || o[AllOps]
}

func (o Operations) Add(op string) {
	o[op] = true
}

func (o Operations) AddAll(operations Operations) {
	for op := range operations {
		o.Add(op)
	}
}

func (o Operations) Remove(op string) {
	delete(o, op)
}

func (o Operations) RemoveAll(operations Operations) {
	for op := range operations {
		o.Remove(op)
	}
}

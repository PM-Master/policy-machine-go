package pip

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

func ToOps(ops ...string) (operations Operations) {
	operations = make(map[string]bool)
	for _, op := range ops {
		operations[op] = true
	}
	return operations
}

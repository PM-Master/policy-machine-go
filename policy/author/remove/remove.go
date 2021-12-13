package remove

import "github.com/PM-Master/policy-machine-go/policy"

func UserAttribute(name string) *policy.DeleteNodeStatement {
	return &policy.DeleteNodeStatement{
		Name: name,
	}
}

func User(name string) *policy.DeleteNodeStatement {
	return &policy.DeleteNodeStatement{
		Name: name,
	}
}

func ObjectAttribute(name string) *policy.DeleteNodeStatement {
	return &policy.DeleteNodeStatement{
		Name: name,
	}
}

func Object(name string) *policy.DeleteNodeStatement {
	return &policy.DeleteNodeStatement{
		Name: name,
	}
}

func PolicyClass(name string) *policy.DeleteNodeStatement {
	return &policy.DeleteNodeStatement{
		Name: name,
	}
}

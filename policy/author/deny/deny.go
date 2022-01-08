package deny

import "github.com/PM-Master/policy-machine-go/policy"

type (
	ops struct {
		subject string
	}

	on struct {
		subject    string
		operations []string
	}

	conts struct {
		subject      string
		operations   []string
		intersection bool
	}
)

func User(user string) *ops {
	return &ops{
		subject: user,
	}
}

func UserAttribute(uattr string) *ops {
	return &ops{
		subject: uattr,
	}
}

func (o *ops) Operations(operations ...string) *on {
	return &on{
		subject:    o.subject,
		operations: operations,
	}
}

func (o *on) On() *conts {
	return &conts{
		subject:    o.subject,
		operations: o.operations,
	}
}

func (c *conts) IntersectionOf() *conts {
	return &conts{
		subject:      c.subject,
		operations:   c.operations,
		intersection: true,
	}
}

func (c *conts) Containers(containers ...string) policy.DenyStatement {
	return policy.DenyStatement{
		Subject:      c.subject,
		Operations:   policy.ToOps(c.operations...),
		Intersection: c.intersection,
		Containers:   containers,
	}
}

func Intersection() bool {
	return true
}

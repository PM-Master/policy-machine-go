package grant

import "github.com/PM-Master/policy-machine-go/policy"

type (
	perms struct {
		ua string
	}

	target struct {
		ua          string
		permissions []string
	}
)

func UserAttribute(name string) *perms {
	return &perms{
		ua: name,
	}
}

func (p *perms) Permissions(perms ...string) *target {
	return &target{
		ua:          p.ua,
		permissions: perms,
	}
}

func (t *target) On(target string) *policy.GrantStatement {
	return &policy.GrantStatement{
		Uattr:      t.ua,
		Target:     target,
		Operations: policy.ToOps(t.permissions...),
	}
}

package create

import "github.com/PM-Master/policy-machine-go/policy"

type (
	when struct {
		label string
	}

	performs struct {
		label   string
		subject policy.Subject
	}

	onDo struct {
		label      string
		subject    policy.Subject
		op         string
		args       []string
		containers []string
	}
)

// create.Obligation("label").When(<subject>).Performs(op, args...).On("", "", "").Do()
// create.Obligation("label").When(<subject>).Performs(op, args...).Do()

func Obligation(label string) *when {
	return &when{label: label}
}

func (w *when) When(subject policy.Subject) *performs {
	return &performs{
		label:   w.label,
		subject: subject,
	}
}

func (p *performs) Performs(op string, args ...string) *onDo {
	return &onDo{
		label:   p.label,
		subject: p.subject,
		op:      op,
		args:    args,
	}
}

func (o *onDo) On(containers []string) *onDo {
	return &onDo{
		label:      o.label,
		subject:    o.subject,
		op:         o.op,
		args:       o.args,
		containers: containers,
	}
}

func (o *onDo) Do(stmts ...policy.Statement) policy.ObligationStatement {
	return policy.ObligationStatement{
		Obligation: policy.Obligation{
			User:  "",
			Label: o.label,
			Event: policy.EventPattern{
				Subject: o.subject,
				Operations: []policy.EventOperation{{
					Operation: o.op,
					Args:      o.args,
				}},
				Containers: o.containers,
			},
			Response: policy.ResponsePattern{
				Actions: stmts,
			},
		},
	}
}

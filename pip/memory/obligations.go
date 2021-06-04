package memory

import (
	"github.com/PM-Master/policy-machine-go/ngac"
)

type memobligations struct {
	obligations map[string]ngac.Obligation
}

func NewObligations() ngac.Obligations {
	return memobligations{
		obligations: make(map[string]ngac.Obligation),
	}
}

func (m memobligations) Add(obligation ngac.Obligation) error {
	m.obligations[obligation.Label] = obligation
	return nil
}

func (m memobligations) Remove(label string) error {
	delete(m.obligations, label)
	return nil
}

func (m memobligations) Get(label string) (ngac.Obligation, error) {
	o := m.obligations[label]
	return ngac.Obligation{
		User:     o.User,
		Label:    o.Label,
		Event:    o.Event,
		Response: o.Response,
	}, nil
}

func (m memobligations) All() ([]ngac.Obligation, error) {
	obligations := make([]ngac.Obligation, 0)

	for _, obligation := range m.obligations {
		obligations = append(obligations, obligation)
	}

	return obligations, nil
}

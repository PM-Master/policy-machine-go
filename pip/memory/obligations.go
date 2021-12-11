package memory

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/policy"
)

type memobligations struct {
	obligations map[string]policy.Obligation
}

func NewObligations() policy.Obligations {
	return &memobligations{
		obligations: make(map[string]policy.Obligation),
	}
}

func (m *memobligations) Add(obligation policy.Obligation) error {
	m.obligations[obligation.Label] = obligation
	return nil
}

func (m *memobligations) Remove(label string) error {
	delete(m.obligations, label)
	return nil
}

func (m *memobligations) Get(label string) (policy.Obligation, error) {
	o := m.obligations[label]
	return policy.Obligation{
		User:     o.User,
		Label:    o.Label,
		Event:    o.Event,
		Response: o.Response,
	}, nil
}

func (m *memobligations) All() ([]policy.Obligation, error) {
	obligations := make([]policy.Obligation, 0)

	for _, obligation := range m.obligations {
		obligations = append(obligations, obligation)
	}

	return obligations, nil
}

func (m *memobligations) MarshalJSON() ([]byte, error) {
	obligations := make([]policy.Obligation, 0)
	for _, obligation := range m.obligations {
		obligations = append(obligations, obligation)
	}

	return json.Marshal(obligations)
}

func (m *memobligations) UnmarshalJSON(bytes []byte) error {
	obligationsArr := make([]policy.Obligation, 0)
	if err := json.Unmarshal(bytes, &obligationsArr); err != nil {
		return err
	}

	obligationsMap := make(map[string]policy.Obligation)
	for _, obligation := range obligationsArr {
		obligationsMap[obligation.Label] = obligation
	}

	m.obligations = obligationsMap

	return nil
}

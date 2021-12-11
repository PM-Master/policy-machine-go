package memory

import (
	"github.com/PM-Master/policy-machine-go/policy"
)

type mempip struct {
	graph        policy.Graph
	prohibitions policy.Prohibitions
	obligations  policy.Obligations
}

func NewPolicyStore() policy.Store {
	return mempip{
		graph:        NewGraph(),
		prohibitions: NewProhibitions(),
		obligations:  NewObligations(),
	}
}

func (m mempip) Graph() policy.Graph {
	return m.graph
}

func (m mempip) Prohibitions() policy.Prohibitions {
	return m.prohibitions
}

func (m mempip) Obligations() policy.Obligations {
	return m.obligations
}

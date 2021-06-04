package memory

import "github.com/PM-Master/policy-machine-go/ngac"

type mempip struct {
	graph        ngac.Graph
	prohibitions ngac.Prohibitions
	obligations  ngac.Obligations
}

func NewPIP() ngac.FunctionalEntity {
	return mempip{
		graph:        NewGraph(),
		prohibitions: NewProhibitions(),
		obligations:  NewObligations(),
	}
}

func (m mempip) Graph() ngac.Graph {
	return m.graph
}

func (m mempip) Prohibitions() ngac.Prohibitions {
	return m.prohibitions
}

func (m mempip) Obligations() ngac.Obligations {
	return m.obligations
}

package memory

import (
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestName(t *testing.T) {
	o := &ngac.Obligation{
		User:  "test",
		Label: "test",
		Event: ngac.EventPattern{
			Subject:    "ANY_USER",
			Operations: []ngac.EventOperation{{"test", []string{"arg1"}}},
			Containers: []string{"!oa1", "oa2"},
		},
		Response: ngac.ResponsePattern{
			Actions: []ngac.Statement{
				&ngac.CreateNodeStatement{
					Name:    "testOA",
					Kind:    graph.ObjectAttribute,
					Parents: []string{"pc1"},
				},
				&ngac.CreatePolicyStatement{
					Name: "testOA",
				},
				&ngac.ObligationStatement{Obligation: ngac.Obligation{
					User:  "bob",
					Label: "label",
					Event: ngac.EventPattern{
						Subject:    "subject",
						Operations: []ngac.EventOperation{{Operation: "op"}},
						Containers: []string{"!oa3"},
					},
					Response: ngac.ResponsePattern{
						Actions: make([]ngac.Statement, 0),
					},
				}},
			},
		},
	}
	bytes, err := o.MarshalJSON()
	require.NoError(t, err)

	o2 := &ngac.Obligation{}
	err = o2.UnmarshalJSON(bytes)
	require.NoError(t, err)
	require.Equal(t, o2, o)
}

func TestJson(t *testing.T) {
	obligations := NewObligations()
	err := obligations.Add(ngac.Obligation{
		User:  "test",
		Label: "test",
		Event: ngac.EventPattern{
			Subject:    "ANY_USER",
			Operations: []ngac.EventOperation{{"test", []string{"arg1"}}},
			Containers: []string{"!oa1", "oa2"},
		},
		Response: ngac.ResponsePattern{
			Actions: []ngac.Statement{
				&ngac.CreateNodeStatement{
					Name:       "testOA",
					Kind:       graph.ObjectAttribute,
					Properties: nil,
					Parents:    []string{"pc1"},
				},
			},
		},
	})
	require.NoError(t, err)

	bytes, err := obligations.MarshalJSON()
	require.NoError(t, err)

	obligations = NewObligations()
	err = obligations.UnmarshalJSON(bytes)
	require.NoError(t, err)
}

package memory

import (
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshal(t *testing.T) {
	o := policy.Obligation{
		User:  "test",
		Label: "test",
		Event: policy.EventPattern{
			Subject:    "ANY_USER",
			Operations: []policy.EventOperation{{"test", []string{"arg1"}}},
			Containers: []string{"!oa1", "oa2"},
		},
		Response: policy.ResponsePattern{
			Actions: []policy.Statement{
				policy.CreateNodeStatement{
					Name:    "testOA",
					Kind:    policy.ObjectAttribute,
					Parents: []string{"pc1"},
				},
				policy.CreatePolicyStatement{
					Name: "testOA",
				},
				policy.ObligationStatement{Obligation: policy.Obligation{
					User:  "bob",
					Label: "label",
					Event: policy.EventPattern{
						Subject:    "subject",
						Operations: []policy.EventOperation{{Operation: "op"}},
						Containers: []string{"!oa3"},
					},
					Response: policy.ResponsePattern{
						Actions: make([]policy.Statement, 0),
					},
				}},
			},
		},
	}
	bytes, err := o.MarshalJSON()
	require.NoError(t, err)

	o2 := policy.Obligation{}
	err = o2.UnmarshalJSON(bytes)
	require.NoError(t, err)
	require.Equal(t, o2, o)
}

func TestJson(t *testing.T) {
	obligations := NewObligations()
	err := obligations.Add(policy.Obligation{
		User:  "test",
		Label: "test",
		Event: policy.EventPattern{
			Subject:    "ANY_USER",
			Operations: []policy.EventOperation{{"test", []string{"arg1"}}},
			Containers: []string{"!oa1", "oa2"},
		},
		Response: policy.ResponsePattern{
			Actions: []policy.Statement{
				policy.CreateNodeStatement{
					Name:       "testOA",
					Kind:       policy.ObjectAttribute,
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

	bytes2, err := obligations.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, bytes, bytes2)
}

package epp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMatches(t *testing.T) {
	evtCtx := EventContext{
		User:   "u1",
		Event:  "read",
		Target: "oa1",
		Args:   map[string]string{},
	}

	pattern := policy.EventPattern{
		Subject:    "ANY_USER",
		Operations: []policy.EventOperation{{Operation: "read"}},
		Containers: []string{},
	}

	matches, err := evtCtx.Matches(pattern)
	require.NoError(t, err)
	require.True(t, matches)

	evtCtx.Event = "write"
	matches, err = evtCtx.Matches(pattern)
	require.NoError(t, err)
	require.False(t, matches)
}

func TestResolveArgs(t *testing.T) {
	args := map[string]string{
		"arg1": "test1",
		"arg2": "test2",
		"arg3": "test3",
		"arg4": "test4",
	}

	t.Run("test create policy", func(t *testing.T) {
		stmt := policy.CreatePolicyStatement{
			Name: "<arg1>",
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.CreatePolicyStatement)
		require.Equal(t, "test1", actual.Name)
	})

	t.Run("test create node", func(t *testing.T) {
		stmt := policy.CreateNodeStatement{
			Name:       "<arg2>",
			Kind:       policy.UserAttribute,
			Properties: map[string]string{"k": "<arg3>"},
			Parents:    []string{"<arg1>", "<arg4>"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.CreateNodeStatement)
		require.Equal(t, "test2", actual.Name)
		require.Equal(t, "test3", actual.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, actual.Parents)
	})

	t.Run("test assign", func(t *testing.T) {
		stmt := policy.AssignStatement{
			Child:   "<arg1>",
			Parents: []string{"<arg2>", "<arg4>"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.AssignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test deassign", func(t *testing.T) {
		stmt := policy.DeassignStatement{
			Child:   "<arg1>",
			Parents: []string{"<arg2>", "<arg4>"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.DeassignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test delete node", func(t *testing.T) {
		stmt := policy.DeleteNodeStatement{
			Name: "<arg1>",
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.DeleteNodeStatement)
		require.Equal(t, "test1", actual.Name)
	})

	t.Run("test grant", func(t *testing.T) {
		stmt := policy.GrantStatement{
			Uattr:  "<arg1>",
			Target: "<arg2>",
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.GrantStatement)
		require.Equal(t, "test1", actual.Uattr)
		require.Equal(t, "test2", actual.Target)
	})

	t.Run("test deny", func(t *testing.T) {
		stmt := policy.DenyStatement{
			Subject:      "<arg1>",
			Operations:   nil,
			Intersection: false,
			Containers:   []string{"!<arg2>", "<arg3>"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.DenyStatement)
		require.Equal(t, "test1", actual.Subject)
		require.Equal(t, []string{"!test2", "test3"}, actual.Containers)
	})

	t.Run("test obligations", func(t *testing.T) {
		stmt := policy.ObligationStatement{
			Obligation: policy.Obligation{
				User:  "",
				Label: "myObl_<arg1>",
				Event: policy.EventPattern{
					Subject:    policy.AnyUserSubject,
					Operations: []policy.EventOperation{{"op1", nil}},
					Containers: []string{"oa1"},
				},
				Response: policy.ResponsePattern{
					Actions: []policy.Statement{
						policy.CreateNodeStatement{
							Name:       "<arg2>",
							Kind:       policy.UserAttribute,
							Properties: map[string]string{"k": "<arg3>"},
							Parents:    []string{"<arg1>", "<arg4>"},
						},
					},
				},
			},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(policy.ObligationStatement)
		require.Equal(t, "myObl_test1", actual.Obligation.Label)

		createNodeStmt := actual.Obligation.Response.Actions[0].(policy.CreateNodeStatement)
		require.Equal(t, "test2", createNodeStmt.Name)
		require.Equal(t, "test3", createNodeStmt.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, createNodeStmt.Parents)
	})
}

func TestSameEventTwice(t *testing.T) {
	policyStore := memory.NewPolicyStore()

	err := policyStore.Graph().CreatePolicyClass("pc1")
	require.NoError(t, err)

	obligation := policy.Obligation{
		User:  "",
		Label: "test_obl",
		Event: policy.EventPattern{
			Subject:    policy.AnyUserSubject,
			Operations: []policy.EventOperation{{Operation: "test_event", Args: []string{"test_arg"}}},
		},
		Response: policy.ResponsePattern{
			Actions: []policy.Statement{
				policy.CreateNodeStatement{
					Name:       "<test_arg>",
					Kind:       policy.ObjectAttribute,
					Properties: nil,
					Parents:    []string{"pc1"},
				},
			},
		},
	}

	err = policyStore.Obligations().Add(obligation)
	require.NoError(t, err)

	epp := NewEPP(policyStore)
	err = epp.ProcessEvent(EventContext{
		User:   "",
		Event:  "test_event",
		Target: "",
		Args:   map[string]string{"test_arg": "foo"},
	})
	require.NoError(t, err)
	err = epp.ProcessEvent(EventContext{
		User:   "",
		Event:  "test_event",
		Target: "",
		Args:   map[string]string{"test_arg": "bar"},
	})
	require.NoError(t, err)

	nodes, err := policyStore.Graph().GetNodes()
	require.NoError(t, err)
	require.Equal(t, 3, len(nodes))
}

func TestName(t *testing.T) {
	c := policy.CreateNodeStatement{}

	func(statement policy.Statement) {
		_, ok := statement.(policy.CreateNodeStatement)
		fmt.Println(ok)
	}(c)
}

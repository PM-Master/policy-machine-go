package epp

import (
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

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.CreatePolicyStatement)
		require.Equal(t, "test1", actual.Name)
	})

	t.Run("test create node", func(t *testing.T) {
		stmt := policy.CreateNodeStatement{
			Name:       "<arg2>",
			Kind:       policy.UserAttribute,
			Properties: map[string]string{"k": "<arg3>"},
			Parents:    []string{"<arg1>", "<arg4>"},
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.CreateNodeStatement)
		require.Equal(t, "test2", actual.Name)
		require.Equal(t, "test3", actual.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, actual.Parents)
	})

	t.Run("test assign", func(t *testing.T) {
		stmt := policy.AssignStatement{
			Child:   "<arg1>",
			Parents: []string{"<arg2>", "<arg4>"},
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.AssignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test deassign", func(t *testing.T) {
		stmt := policy.DeassignStatement{
			Child:   "<arg1>",
			Parents: []string{"<arg2>", "<arg4>"},
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.DeassignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test delete node", func(t *testing.T) {
		stmt := policy.DeleteNodeStatement{
			Name: "<arg1>",
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.DeleteNodeStatement)
		require.Equal(t, "test1", actual.Name)
	})

	t.Run("test grant", func(t *testing.T) {
		stmt := policy.GrantStatement{
			Uattr:  "<arg1>",
			Target: "<arg2>",
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.GrantStatement)
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

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.DenyStatement)
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
						&policy.CreateNodeStatement{
							Name:       "<arg2>",
							Kind:       policy.UserAttribute,
							Properties: map[string]string{"k": "<arg3>"},
							Parents:    []string{"<arg1>", "<arg4>"},
						},
					},
				},
			},
		}

		resolved, err := resolveArgs(&stmt, args)
		require.NoError(t, err)
		actual := resolved.(*policy.ObligationStatement)
		require.Equal(t, "myObl_test1", actual.Obligation.Label)

		createNodeStmt := actual.Obligation.Response.Actions[0].(*policy.CreateNodeStatement)
		require.Equal(t, "test2", createNodeStmt.Name)
		require.Equal(t, "test3", createNodeStmt.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, createNodeStmt.Parents)
	})
}

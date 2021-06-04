package epp

import (
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
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

	pattern := ngac.EventPattern{
		Subject:    "ANY_USER",
		Operations: []ngac.EventOperation{{Operation: "read"}},
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
		"$arg1": "test1",
		"$arg2": "test2",
		"$arg3": "test3",
		"$arg4": "test4",
	}

	t.Run("test create policy", func(t *testing.T) {
		stmt := ngac.CreatePolicyStatement{
			Name: "$arg1",
			Statements: []ngac.Statement{
				ngac.CreateNodeStatement{
					Name:       "$arg2",
					Kind:       graph.UserAttribute,
					Properties: map[string]string{"k": "$arg3"},
					Parents:    []string{"$arg4"},
				},
			},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.CreatePolicyStatement)
		require.Equal(t, "test1", actual.Name)

		createNodeStmt := actual.Statements[0].(ngac.CreateNodeStatement)
		require.Equal(t, "test2", createNodeStmt.Name)
		require.Equal(t, "test3", createNodeStmt.Properties["k"])
		require.Equal(t, []string{"test4"}, createNodeStmt.Parents)
	})

	t.Run("test create node", func(t *testing.T) {
		stmt := ngac.CreateNodeStatement{
			Name:       "$arg2",
			Kind:       graph.UserAttribute,
			Properties: map[string]string{"k": "$arg3"},
			Parents:    []string{"$arg1", "$arg4"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.CreateNodeStatement)
		require.Equal(t, "test2", actual.Name)
		require.Equal(t, "test3", actual.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, actual.Parents)
	})

	t.Run("test assign", func(t *testing.T) {
		stmt := ngac.AssignStatement{
			Child:   "$arg1",
			Parents: []string{"$arg2", "$arg4"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.AssignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test deassign", func(t *testing.T) {
		stmt := ngac.DeassignStatement{
			Child:   "$arg1",
			Parents: []string{"$arg2", "$arg4"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.DeassignStatement)
		require.Equal(t, "test1", actual.Child)
		require.Equal(t, []string{"test2", "test4"}, actual.Parents)
	})

	t.Run("test delete node", func(t *testing.T) {
		stmt := ngac.DeleteNodeStatement{
			Name: "$arg1",
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.DeleteNodeStatement)
		require.Equal(t, "test1", actual.Name)
	})

	t.Run("test grant", func(t *testing.T) {
		stmt := ngac.GrantStatement{
			Uattr:  "$arg1",
			Target: "$arg2",
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.GrantStatement)
		require.Equal(t, "test1", actual.Uattr)
		require.Equal(t, "test2", actual.Target)
	})

	t.Run("test deny", func(t *testing.T) {
		stmt := ngac.DenyStatement{
			Subject:      "$arg1",
			Operations:   nil,
			Intersection: false,
			Containers:   []string{"!$arg2", "$arg3"},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.DenyStatement)
		require.Equal(t, "test1", actual.Subject)
		require.Equal(t, []string{"!test2", "test3"}, actual.Containers)
	})

	t.Run("test obligations", func(t *testing.T) {
		stmt := ngac.ObligationStatement{
			Obligation: ngac.Obligation{
				User:  "",
				Label: "myObl_$arg1",
				Event: ngac.EventPattern{
					Subject:    "ANY_USER",
					Operations: []ngac.EventOperation{{"op1", nil}},
					Containers: []string{"oa1"},
				},
				Response: ngac.ResponsePattern{
					Actions: []ngac.Statement{
						ngac.CreateNodeStatement{
							Name:       "$arg2",
							Kind:       graph.UserAttribute,
							Properties: map[string]string{"k": "$arg3"},
							Parents:    []string{"$arg1", "$arg4"},
						},
					},
				},
			},
		}

		resolved, err := resolveArgs(stmt, args)
		require.NoError(t, err)
		actual := resolved.(ngac.ObligationStatement)
		require.Equal(t, "myObl_test1", actual.Obligation.Label)

		createNodeStmt := actual.Obligation.Response.Actions[0].(ngac.CreateNodeStatement)
		require.Equal(t, "test2", createNodeStmt.Name)
		require.Equal(t, "test3", createNodeStmt.Properties["k"])
		require.Equal(t, []string{"test1", "test4"}, createNodeStmt.Parents)
	})
}

package author

import (
	"github.com/PM-Master/policy-machine-go/epp"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/PM-Master/policy-machine-go/policy/author/assign"
	"github.com/PM-Master/policy-machine-go/policy/author/create"
	"github.com/PM-Master/policy-machine-go/policy/author/deassign"
	"github.com/PM-Master/policy-machine-go/policy/author/deny"
	"github.com/PM-Master/policy-machine-go/policy/author/grant"
	"github.com/stretchr/testify/require"
	"testing"
)

/**
policy.Author
author.Author
policy.Author(store, s, s, s, s, s)
author.Author(store, s, s, s, s, s)



*/

func TestNew(t *testing.T) {
	assign.UserAttribute("ua1").To("pc1")
	create.UserAttribute("ua1").In("pc1")
	create.ObjectAttribute("ua1").In("pc1")
	deny.User("bob").
		Operations("read").
		On().
		IntersectionOf().
		Containers("oa1", "!oa2")

	grant.UserAttribute("ua1").
		Permissions("read", "write").
		On("target")

	create.Obligation("test").
		When("subject").
		Performs("op", "arg1", "arg2").
		Do(
			create.
				Object("test_object").
				In("oa1"),
		)
}

func TestAuthor(t *testing.T) {
	pip := memory.NewPolicyStore()
	err := Author(pip,
		create.PolicyClass("pc1"),
		create.UserAttribute("ua1").In("pc1"),
		create.UserAttribute("ua2").WithProperties("k", "v").In("pc1"),
		create.ObjectAttribute("oa1").In("pc1"),
		create.User("u1").In("ua1"),
		create.Object("o1").In("oa1"),

		assign.User("u1").To("ua2"),
		deassign.User("u1").From("ua1"),

		grant.UserAttribute("ua1").Permissions("read", "write").On("oa1"),

		deny.User("u1").Operations("write").On().IntersectionOf().Containers("oa1", "!oa1"),

		create.Obligation("obligation1").When(policy.AnyUserSubject).Performs("test_op", "foo", "bar").Do(
			create.UserAttribute("<foo>").In("ua2"),
			create.UserAttribute("<bar>").In("ua2"),
		),
	)
	require.NoError(t, err)

	evtProc := epp.NewEPP(pip)
	err = evtProc.ProcessEvent(epp.EventContext{
		User:   "",
		Event:  "test_op",
		Target: "",
		Args:   map[string]string{"foo": "hello", "bar": "world"},
	})
	require.NoError(t, err)

	exists, err := pip.Graph().Exists("hello")
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = pip.Graph().Exists("world")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestApply(t *testing.T) {
	pip := memory.NewPolicyStore()
	err := ReadAndApply(pip, "testdata/test.ngac")
	require.NoError(t, err)

	nodes, err := pip.Graph().GetNodes()
	require.NoError(t, err)
	require.Contains(t, nodes, "rbac")
	require.Contains(t, nodes, "ua1")
	require.Contains(t, nodes, "ua2")
	require.Contains(t, nodes, "oa1")
	require.Contains(t, nodes, "oa2")

	parents, err := pip.Graph().GetParents("ua1")
	require.NoError(t, err)
	require.Contains(t, parents, "rbac")

	parents, err = pip.Graph().GetParents("ua2")
	require.NoError(t, err)
	require.Contains(t, parents, "rbac")

	parents, err = pip.Graph().GetParents("oa1")
	require.NoError(t, err)
	require.Contains(t, parents, "rbac")

	parents, err = pip.Graph().GetParents("oa2")
	require.NoError(t, err)
	require.Contains(t, parents, "rbac")
}

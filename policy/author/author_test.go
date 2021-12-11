package author

import (
	"github.com/PM-Master/policy-machine-go/epp"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/stretchr/testify/require"
	"testing"
)

/**
policy.Author
author.Author
policy.Author(store, s, s, s, s, s)
author.Author(store, s, s, s, s, s)
*/

func TestAuthor(t *testing.T) {
	pip := memory.NewPolicyStore()
	err := Author(pip,
		PolicyClass("pc1"),
		UserAttribute("ua1", "pc1"),
		UserAttributeWithProperties("ua2", Properties{"k": "v"}, "pc1"),

		ObjectAttribute("oa1", "pc1"),
		ObjectAttributeWithProperties("oa2", Properties{"k": "v"}, "pc1"),

		User("u1", "ua1"),
		UserWithProperties("u2", Properties{"k": "v"}, "ua1"),

		Object("o1", "oa1"),
		ObjectWithProperties("o2", Properties{"k": "v"}, "oa1"),

		Assign("u1", "ua2"),

		Deassign("u1", "ua1"),

		Grant("ua1", "oa1", "read", "write"),

		Deny("u1", policy.ToOps("write"), false, "oa1", "!oa1"),

		Obligation("obligation1",
			Event(policy.AnyUserSubject, Operation("test_op", "foo", "bar")),
			Response(
				UserAttribute("<foo>", "ua2"),
				UserAttribute("<bar>", "ua2"),
			),
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

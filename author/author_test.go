package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestApply(t *testing.T) {
	pip := memory.NewPIP()
	author := New(pip)
	fmt.Println(os.Getwd())
	err := author.ReadAndApply("testdata/test.ngac")
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

func TestGraphAndObligations(t *testing.T) {
	pip := memory.NewPIP()
	author := New(pip)
	err := author.ReadAndApply("testdata/test2.ngac")
	require.NoError(t, err)
}

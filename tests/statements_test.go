package tests

import (
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePolicyClassStatement(t *testing.T) {
	stmt := ngac.CreatePolicyStatement{
		Name: "testpc",
	}

	pip := memory.NewPIP()
	err := stmt.Apply(pip)
	require.NoError(t, err)
	exists, err := pip.Graph().Exists("testpc")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestCreatePolicyClassStatementWithStatements(t *testing.T) {
	stmt := ngac.CreatePolicyStatement{
		Name: "testpc",
	}

	pip := memory.NewPIP()
	err := stmt.Apply(pip)
	require.NoError(t, err)
	exists, err := pip.Graph().Exists("testpc")
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = pip.Graph().Exists("ua1")
	require.NoError(t, err)
	require.True(t, exists)
	parents, err := pip.Graph().GetParents("ua1")
	require.NoError(t, err)
	require.Contains(t, parents, "testpc")
}

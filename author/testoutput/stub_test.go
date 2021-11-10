package testoutput

import (
	"github.com/PM-Master/policy-machine-go/author"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStubFunctionCall(t *testing.T) {
	pip := memory.NewPIP()
	author := author.New(pip)
	err := author.ReadPAL("test.ngac")
	require.NoError(t, err)
	err = pip.Graph().CreatePolicyClass("rbac")
	require.NoError(t, err)
	require.NoError(t, err)
	stub := NewTestStub(author, pip)
	err = stub.my_function1("myarg1", "myarg2")
	require.NoError(t, err)
	exists, err := pip.Graph().Exists("myarg1")
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = pip.Graph().Exists("myarg2_test")
	require.NoError(t, err)
	require.True(t, exists)
}

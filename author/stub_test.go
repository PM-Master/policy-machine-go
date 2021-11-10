package author

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateFunctionsStub(t *testing.T) {
	err := GenerateFunctionsStub("TestStub", "testdata", "testoutput/stub.go")
	require.NoError(t, err)
}

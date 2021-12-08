package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	obligationStr := `
OBLIGATION myObl
WHEN ANY_USER
PERFORMS op1(arg1,arg2)
ON oa1
DO(
    create user attribute $arg1 in ua2;
    create user attribute $arg2 in ua4;

	OBLIGATION myObl
	WHEN ANY_USER
	PERFORMS op2 OR op3
	DO(
		assign ua3 to ua2;
	);
);
`

	parser := NewObligationParser()
	obligation, err := parser.Parse(obligationStr)
	require.NoError(t, err)

	require.Equal(t, "myObl", obligation.Label)
	require.Equal(t, "ANY_USER", obligation.Event.Subject)
	require.Equal(t, []ngac.EventOperation{{
		Operation: "op1",
		Args:      []string{"arg1", "arg2"},
	}}, obligation.Event.Operations)
	require.Equal(t, 3, len(obligation.Response.Actions))
}

func TestParserNoOnClause(t *testing.T) {
	obligationStr := `
OBLIGATION myObl
WHEN ANY_USER
PERFORMS op1(arg1,arg2)
DO(
    create user attribute $arg1 in ua2;
    create user attribute $arg2 in ua4;

	OBLIGATION myObl
	WHEN ANY_USER
	PERFORMS op2 OR op3
	DO(
		assign ua3 to ua2;
	);
);
`

	parser := NewObligationParser()
	obligation, err := parser.Parse(obligationStr)
	require.NoError(t, err)

	require.Equal(t, "myObl", obligation.Label)
	require.Equal(t, "ANY_USER", obligation.Event.Subject)
	require.Equal(t, []ngac.EventOperation{{
		Operation: "op1",
		Args:      []string{"arg1", "arg2"},
	}}, obligation.Event.Operations)
	require.Equal(t, 3, len(obligation.Response.Actions))
}

func TestName(t *testing.T) {
	strings := []string{"1", "2", "3"}
	fmt.Println(strings[2:3])
}

package author

import (
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"github.com/stretchr/testify/require"
	"testing"
)

var testStr = `
create policy RBAC(
    create user attribute ua1;
);

create user bob IN ua3;
create resource resource1 IN oa1;

OBLIGATION obl_label
WHEN ANY_USER
PERFORMS op1(arg1, arg2)
DO (
    create user attribute ua4 IN ua1; 

    OBLIGATION obl_label
    WHEN ANY_USER
    PERFORMS op2 OR op3
    DO (
        create user attribute ua4 IN ua1;
    );
);
`

func TestParse(t *testing.T) {
	Parse(testStr)
}

func TestParseCreatePolicyClass(t *testing.T) {
	s := `create policy RBAC(
    		create user attribute ua1 with properties k1=v1, k2=v2 IN ua2, ua3;
		);`
	stmt, err := parseCreatePolicy(s)
	require.NoError(t, err)

	createPolicyStmt := stmt.(ngac.CreatePolicyStatement)
	require.Equal(t, "RBAC", createPolicyStmt.Name)
	require.Equal(t, 1, len(createPolicyStmt.Statements))
	stmt = createPolicyStmt.Statements[0]
	createUAStmt := stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "ua1", createUAStmt.Name)
	require.Equal(t, graph.UserAttribute, createUAStmt.Kind)
	require.Equal(t, map[string]string{"k1": "v1", "k2": "v2"}, createUAStmt.Properties)
	require.Equal(t, []string{"ua2", "ua3"}, createUAStmt.Parents)
}

func TestParseCreateNode(t *testing.T) {
	s := "create user u1 in ua1"
	stmt, err := parseCreateNode(s)
	require.NoError(t, err)
	nodeStmt := stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "u1", nodeStmt.Name)
	require.Equal(t, graph.User, nodeStmt.Kind)
	require.Equal(t, []string{"ua1"}, nodeStmt.Parents)

	s = "create object o1 in oa1"
	stmt, err = parseCreateNode(s)
	require.NoError(t, err)
	nodeStmt = stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "o1", nodeStmt.Name)
	require.Equal(t, graph.Object, nodeStmt.Kind)
	require.Equal(t, []string{"oa1"}, nodeStmt.Parents)

	s = "create user attribute ua1 in ua2"
	stmt, err = parseCreateNode(s)
	require.NoError(t, err)
	nodeStmt = stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "ua1", nodeStmt.Name)
	require.Equal(t, graph.UserAttribute, nodeStmt.Kind)
	require.Equal(t, []string{"ua2"}, nodeStmt.Parents)

	s = "create object attribute oa1 in oa2"
	stmt, err = parseCreateNode(s)
	require.NoError(t, err)
	nodeStmt = stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "oa1", nodeStmt.Name)
	require.Equal(t, graph.ObjectAttribute, nodeStmt.Kind)
	require.Equal(t, []string{"oa2"}, nodeStmt.Parents)

	s = "create object attribute oa1 with properties k1=v1, k2=v2 in oa2"
	stmt, err = parseCreateNode(s)
	require.NoError(t, err)
	nodeStmt = stmt.(ngac.CreateNodeStatement)
	require.Equal(t, "oa1", nodeStmt.Name)
	require.Equal(t, graph.ObjectAttribute, nodeStmt.Kind)
	require.Equal(t, []string{"oa2"}, nodeStmt.Parents)
	require.Equal(t, map[string]string{"k1": "v1", "k2": "v2"}, nodeStmt.Properties)
}

func TestParseDeny(t *testing.T) {
	s := "deny ua1 read, write on !oa1, oa2"
	stmt, err := parseDeny(s)
	require.NoError(t, err)
	denyStmt := stmt.(ngac.DenyStatement)
	require.Equal(t, "ua1", denyStmt.Subject)
	require.Equal(t, graph.ToOps("read", "write"), denyStmt.Operations)
	require.Equal(t, false, denyStmt.Intersection)
	require.Equal(t, []string{"!oa1", "oa2"}, denyStmt.Containers)

	s = "deny ua1 read, write on intersection of !oa1, oa2"
	stmt, err = parseDeny(s)
	require.NoError(t, err)
	denyStmt = stmt.(ngac.DenyStatement)
	require.Equal(t, "ua1", denyStmt.Subject)
	require.Equal(t, graph.ToOps("read", "write"), denyStmt.Operations)
	require.Equal(t, true, denyStmt.Intersection)
	require.Equal(t, []string{"!oa1", "oa2"}, denyStmt.Containers)
}

func TestParseGrant(t *testing.T) {
	s := "grant ua1 read, write on oa2"
	stmt, err := parseGrant(s)
	require.NoError(t, err)
	grantStmt := stmt.(ngac.GrantStatement)
	require.Equal(t, "ua1", grantStmt.Uattr)
	require.Equal(t, "oa2", grantStmt.Target)
	require.Equal(t, graph.ToOps("read", "write"), grantStmt.Operations)
}

func TestParseAssign(t *testing.T) {
	s := "assign ua1 to ua2, ua3"
	stmt, err := parseAssign(s)
	require.NoError(t, err)
	assignStmt := stmt.(ngac.AssignStatement)
	require.Equal(t, "ua1", assignStmt.Child)
	require.Equal(t, []string{"ua2", "ua3"}, assignStmt.Parents)
}

func TestParseDeassign(t *testing.T) {
	s := "deassign ua1 FROM ua2, ua3"
	stmt, err := parseDeassign(s)
	require.NoError(t, err)
	deassignStmt := stmt.(ngac.DeassignStatement)
	require.Equal(t, "ua1", deassignStmt.Child)
	require.Equal(t, []string{"ua2", "ua3"}, deassignStmt.Parents)
}

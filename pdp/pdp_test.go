package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
	"testing"
)

func TestPDP(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1")
	g.CreateNode("oa1", policy.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", policy.UserAttribute, nil, "pc1")
	g.CreateNode("o1", policy.Object, nil, "oa1")
	g.CreateNode("u1", policy.User, nil, "ua1")
	g.Associate("ua1", "oa1", policy.ToOps("r", "w"))

	decider := NewDecider(g, nil)
	actual, _ := decider.HasPermissions("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

func TestPDP2(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1")
	g.CreatePolicyClass("pc2")
	g.CreateNode("oa1", policy.ObjectAttribute, nil, "pc1")
	g.CreateNode("oa2", policy.ObjectAttribute, nil, "pc2")
	g.CreateNode("ua1", policy.UserAttribute, nil, "pc1")
	g.CreateNode("o1", policy.Object, nil, "oa1", "oa2")
	g.CreateNode("u1", policy.User, nil, "ua1")
	g.Associate("ua1", "oa1", policy.ToOps("r", "w"))
	g.Associate("ua1", "oa2", policy.ToOps("r"))

	decider := NewDecider(g, nil)
	actual, _ := decider.HasPermissions("u1", "o1", "r", "w")
	if actual {
		t.Fatal("u1 should not have [w] on o1")
	}

	actual, _ = decider.HasPermissions("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

func TestPDPAllOps(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1")
	g.CreateNode("oa1", policy.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", policy.UserAttribute, nil, "pc1")
	g.CreateNode("o1", policy.Object, nil, "oa1")
	g.CreateNode("u1", policy.User, nil, "ua1")
	g.Associate("ua1", "oa1", policy.ToOps("*"))

	decider := NewDecider(g, nil)
	actual, _ := decider.HasPermissions("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

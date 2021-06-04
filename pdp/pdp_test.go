package pdp

import (
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"testing"
)

func TestPDP(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1")
	g.CreateNode("oa1", graph.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", graph.UserAttribute, nil, "pc1")
	g.CreateNode("o1", graph.Object, nil, "oa1")
	g.CreateNode("u1", graph.User, nil, "ua1")
	g.Associate("ua1", "oa1", graph.ToOps("r", "w"))

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
	g.CreateNode("oa1", graph.ObjectAttribute, nil, "pc1")
	g.CreateNode("oa2", graph.ObjectAttribute, nil, "pc2")
	g.CreateNode("ua1", graph.UserAttribute, nil, "pc1")
	g.CreateNode("o1", graph.Object, nil, "oa1", "oa2")
	g.CreateNode("u1", graph.User, nil, "ua1")
	g.Associate("ua1", "oa1", graph.ToOps("r", "w"))
	g.Associate("ua1", "oa2", graph.ToOps("r"))

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
	g.CreateNode("oa1", graph.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", graph.UserAttribute, nil, "pc1")
	g.CreateNode("o1", graph.Object, nil, "oa1")
	g.CreateNode("u1", graph.User, nil, "ua1")
	g.Associate("ua1", "oa1", graph.ToOps("*"))

	decider := NewDecider(g, nil)
	actual, _ := decider.HasPermissions("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

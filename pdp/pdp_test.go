package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"testing"
)

func TestPDP(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1", nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", pip.UserAttribute, nil, "pc1")
	g.CreateNode("o1", pip.Object, nil, "oa1")
	g.CreateNode("u1", pip.User, nil, "ua1")
	g.Associate("ua1", "oa1", pip.ToOps("r", "w"))

	decider := NewDecider(g)
	actual := decider.Decide("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

func TestPDP2(t *testing.T) {
	g := memory.NewGraph()
	g.CreatePolicyClass("pc1", nil)
	g.CreatePolicyClass("pc2", nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil, "pc1")
	g.CreateNode("oa2", pip.ObjectAttribute, nil, "pc2")
	g.CreateNode("ua1", pip.UserAttribute, nil, "pc1")
	g.CreateNode("o1", pip.Object, nil, "oa1", "oa2")
	g.CreateNode("u1", pip.User, nil, "ua1")
	g.Associate("ua1", "oa1", pip.ToOps("r", "w"))
	g.Associate("ua1", "oa2", pip.ToOps("r"))

	decider := NewDecider(g)
	actual := decider.Decide("u1", "o1", "r", "w")
	if actual {
		t.Fatal("u1 should not have [w] on o1")
	}

	actual = decider.Decide("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

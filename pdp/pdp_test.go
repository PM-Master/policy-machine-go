package pdp

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"testing"
)

func TestPDP(t *testing.T) {
	g := memory.NewGraph()
	g.CreateNode("pc1", pip.PolicyClass, nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil)
	g.CreateNode("ua1", pip.UserAttribute, nil)
	g.CreateNode("o1", pip.Object, nil)
	g.CreateNode("u1", pip.User, nil)
	g.Assign("u1", "ua1")
	g.Assign("o1", "oa1")
	g.Assign("oa1", "pc1")
	g.Assign("ua1", "pc1")
	g.Associate("ua1", "oa1", pip.ToOps("r", "w"))

	decider := NewDecider(g)
	actual, _ := decider.Decide("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

func TestPDP2(t *testing.T) {
	g := memory.NewGraph()
	g.CreateNode("pc1", pip.PolicyClass, nil)
	g.CreateNode("pc2", pip.PolicyClass, nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil)
	g.CreateNode("oa2", pip.ObjectAttribute, nil)
	g.CreateNode("ua1", pip.UserAttribute, nil)
	g.CreateNode("o1", pip.Object, nil)
	g.CreateNode("u1", pip.User, nil)
	g.Assign("u1", "ua1")
	g.Assign("o1", "oa1")
	g.Assign("o1", "oa2")
	g.Assign("oa1", "pc1")
	g.Assign("oa2", "pc2")
	g.Assign("ua1", "pc1")
	g.Associate("ua1", "oa1", pip.ToOps("r", "w"))
	g.Associate("ua1", "oa2", pip.ToOps("r"))

	decider := NewDecider(g)
	actual, _ := decider.Decide("u1", "o1", "r", "w")
	if actual {
		t.Fatal("u1 should not have [w] on o1")
	}

	actual, _ = decider.Decide("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

func TestPDPAllOps(t *testing.T) {
	g := memory.NewGraph()
	g.CreateNode("pc1", pip.PolicyClass, nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil)
	g.CreateNode("ua1", pip.UserAttribute, nil)
	g.CreateNode("o1", pip.Object, nil)
	g.CreateNode("u1", pip.User, nil)
	g.Assign("u1", "ua1")
	g.Assign("o1", "oa1")
	g.Assign("oa1", "pc1")
	g.Assign("ua1", "pc1")
	g.Associate("ua1", "oa1", pip.ToOps("*"))

	decider := NewDecider(g)
	actual, _ := decider.Decide("u1", "o1", "r")
	if !actual {
		t.Fatal("u1 should have [r] on o1")
	}
}

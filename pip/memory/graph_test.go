package memory

import (
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"testing"
)

func TestJSON(t *testing.T) {
	g := NewGraph()
	g.CreatePolicyClass("pc1")
	g.CreateNode("oa1", graph.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", graph.UserAttribute, nil, "pc1")
	g.CreateNode("o1", graph.Object, nil, "oa1")
	g.CreateNode("u1", graph.User, nil, "ua1")
	g.Associate("ua1", "oa1", graph.ToOps("r", "w"))

	b, _ := g.MarshalJSON()
	g1 := NewGraph()
	g1.UnmarshalJSON(b)

	if ok, _ := g1.Exists("pc1"); !ok {
		t.Fatal("pc1 should exist but does not")
	}
	if ok, _ := g1.Exists("oa1"); !ok {
		t.Fatal("oa1 should exist but does not")
	}
	if ok, _ := g1.Exists("ua1"); !ok {
		t.Fatal("ua1 should exist but does not")
	}
	if ok, _ := g1.Exists("o1"); !ok {
		t.Fatal("o1 should exist but does not")
	}
	if ok, _ := g1.Exists("u1"); !ok {
		t.Fatal("u1 should exist but does not")
	}
}

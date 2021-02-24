package memory

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"testing"
)

func TestJSON(t *testing.T) {
	g := NewGraph()
	g.CreatePolicyClass("pc1", nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil, "pc1")
	g.CreateNode("ua1", pip.UserAttribute, nil, "pc1")
	g.CreateNode("o1", pip.Object, nil, "oa1")
	g.CreateNode("u1", pip.User, nil, "ua1")
	g.Associate("ua1", "oa1", pip.ToOps("r", "w"))

	b, _ := g.MarshalJSON()
	g.UnmarshalJSON(b)

	if ok, _ := g.Exists("pc1"); !ok {
		t.Fatal("pc1 should exist but does not")
	}
	if ok, _ := g.Exists("oa1"); !ok {
		t.Fatal("oa1 should exist but does not")
	}
	if ok, _ := g.Exists("ua1"); !ok {
		t.Fatal("ua1 should exist but does not")
	}
	if ok, _ := g.Exists("o1"); !ok {
		t.Fatal("o1 should exist but does not")
	}
	if ok, _ := g.Exists("u1"); !ok {
		t.Fatal("u1 should exist but does not")
	}
}

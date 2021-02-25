package memory

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"testing"
)

func TestJSON(t *testing.T) {
	g := NewGraph()
	g.CreateNode("pc1", pip.PolicyClass, nil)
	g.CreateNode("oa1", pip.ObjectAttribute, nil)
	g.CreateNode("ua1", pip.UserAttribute, nil)
	g.CreateNode("o1", pip.Object, nil)
	g.CreateNode("u1", pip.User, nil)
	g.Assign("oa1", "pc1")
	g.Assign("ua1", "pc1")
	g.Assign("u1", "ua1")
	g.Assign("o1", "oa1")
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

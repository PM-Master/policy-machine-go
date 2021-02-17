package pip

import (
	"testing"
)

func TestJSON(t *testing.T) {
	g := NewGraph()
	g.CreateNode("pc1", PolicyClass, nil)
	g.CreateNode("oa1", ObjectAttribute, nil)
	g.CreateNode("ua1", UserAttribute, nil)
	g.CreateNode("o1", Object, nil)
	g.CreateNode("u1", User, nil)
	g.Assign("u1", "ua1")
	g.Assign("o1", "oa1")
	g.Assign("oa1", "pc1")
	g.Assign("ua1", "pc1")
	g.Associate("ua1", "oa1", ToOps("r", "w"))

	str := ToJson(g)
	g1 := FromJson(str)

	if _, ok := g1.GetNode("pc1"); !ok {
		t.Fatal("pc1 should exist but does not")
	}
	if _, ok := g1.GetNode("oa1"); !ok {
		t.Fatal("oa1 should exist but does not")
	}
	if _, ok := g1.GetNode("ua1"); !ok {
		t.Fatal("ua1 should exist but does not")
	}
	if _, ok := g1.GetNode("o1"); !ok {
		t.Fatal("o1 should exist but does not")
	}
	if _, ok := g1.GetNode("u1"); !ok {
		t.Fatal("u1 should exist but does not")
	}
}

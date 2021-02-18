package pip

import (
	"encoding/json"
	"fmt"
)

type (
	Graph struct {
		Nodes        map[string]Node                  `json:"nodes"`
		Assignments  map[string]map[string]bool       `json:"assignments"`
		Associations map[string]map[string]Operations `json:"associations"`
	}
)

func NewGraph() Graph {
	return Graph{
		Nodes:        make(map[string]Node),
		Assignments:  make(map[string]map[string]bool),
		Associations: make(map[string]map[string]Operations),
	}
}

func (g Graph) CreateNode(name string, kind Kind, properties map[string]string) error {
	if _, ok := g.Nodes[name]; ok {
		return fmt.Errorf("name %q is already exists", name)
	}

	if properties == nil {
		properties = make(map[string]string)
	}

	g.Nodes[name] = Node{
		Name:       name,
		Kind:       kind,
		Properties: properties,
	}

	return nil
}

func (g Graph) DeleteNode(name string) {
	delete(g.Nodes, name)
}

func (g Graph) GetNode(name string) (Node, bool) {
	node, ok := g.Nodes[name]
	return node, ok
}

func (g Graph) Find(kind Kind, properties map[string]string) map[string]Node {
	found := make(map[string]Node)
	for _, node := range g.Nodes {
		if node.Kind != kind {
			continue
		}

		match := true
		for k, v := range properties {
			if node.Properties[k] != v {
				match = false
			}
		}

		if match {
			found[node.Name] = node
		}
	}

	return found
}

func (g Graph) Assign(source string, target string) {
	if _, ok := g.Assignments[source]; !ok {
		g.Assignments[source] = make(map[string]bool)
	}
	g.Assignments[source][target] = true
}

func (g Graph) Deassign(source string, target string) {
	delete(g.Assignments[source], target)
}

func (g Graph) GetChildren(name string) map[string]Node {
	children := make(map[string]Node)
	for nodeName, assignmentMap := range g.Assignments {
		if assignmentMap[name] {
			node, _ := g.GetNode(nodeName)
			children[nodeName] = node
		}
	}
	return children
}

func (g Graph) GetParents(name string) map[string]Node {
	assignments := g.Assignments[name]
	parents := make(map[string]Node)
	for nodeName := range assignments {
		node, _ := g.GetNode(nodeName)
		parents[nodeName] = node
	}
	return parents
}

func (g Graph) Associate(source string, target string, operations Operations) {
	if _, ok := g.Associations[source]; !ok {
		g.Associations[source] = make(map[string]Operations)
	}
	g.Associations[source][target] = operations
}

func (g Graph) Dissociate(source string, target string) {
	delete(g.Associations[source], target)
}

func (g Graph) GetAssociations(name string) map[string]Operations {
	return g.Associations[name]
}

func ToJson(graph Graph) string {
	b, _ := json.Marshal(graph)
	return string(b)
}

func FromJson(s string) Graph {
	g := NewGraph()
	if err := json.Unmarshal([]byte(s), &g); err != nil {

	}
	return g
}

package memory

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/pip"
)

type (
	graph struct {
		nodes        map[string]pip.Node
		assignments  map[string]map[string]bool
		associations map[string]map[string]pip.Operations
	}
)

func NewGraph() pip.Graph {
	return &graph{
		nodes:        make(map[string]pip.Node),
		assignments:  make(map[string]map[string]bool),
		associations: make(map[string]map[string]pip.Operations),
	}
}

func (g graph) CreateNode(name string, kind pip.Kind, properties map[string]string) (pip.Node, error) {
	if _, ok := g.nodes[name]; ok {
		return pip.Node{}, fmt.Errorf("name %q is already exists", name)
	}

	if properties == nil {
		properties = make(map[string]string)
	}

	n := pip.Node{
		Name:       name,
		Kind:       kind,
		Properties: properties,
	}
	node := copyNode(n)
	g.nodes[name] = node

	return node, nil
}

func (g graph) UpdateNode(name string, properties map[string]string) error {
	if ok, _ := g.Exists(name); !ok {
		return fmt.Errorf("node %q does not exist", name)
	}

	n := g.nodes[name]
	n.Properties = properties
	g.nodes[name] = copyNode(n)

	return nil
}

func (g graph) DeleteNode(name string) error {
	// delete this node's assignments
	// return an error if this node has other nodes assigned to it still
	if children, _ := g.GetChildren(name); len(children) > 0 {
		return fmt.Errorf("cannot delete %q because it has nodes assigned to it", name)
	}

	delete(g.assignments, name)

	// delete associations
	assocs, _ := g.GetAssociations()
	for subject, subjectAssocs := range assocs {
		if _, ok := subjectAssocs[name]; !ok {
			continue
		}

		delete(g.associations[subject], name)
	}

	delete(g.associations, name)

	// delete node
	delete(g.nodes, name)

	return nil
}

func (g graph) Exists(name string) (bool, error) {
	_, ok := g.nodes[name]
	return ok, nil
}

func (g graph) GetNodes() (map[string]pip.Node, error) {
	nodes := make(map[string]pip.Node)
	for _, node := range g.nodes {
		copyNode := copyNode(node)
		nodes[copyNode.Name] = copyNode
	}
	return nodes, nil
}

func (g graph) GetNode(name string) (pip.Node, error) {
	node, ok := g.nodes[name]
	if !ok {
		return pip.Node{}, fmt.Errorf("node %q does not exist", name)
	}
	return copyNode(node), nil
}

func (g graph) Find(kind pip.Kind, properties map[string]string) (map[string]pip.Node, error) {
	found := make(map[string]pip.Node)
	for _, node := range g.nodes {
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

	return found, nil
}

func (g graph) Assign(child string, parent string) error {
	var (
		childNode  pip.Node
		parentNode pip.Node
		err        error
	)

	if childNode, err = g.GetNode(child); err != nil {
		return err
	}
	if parentNode, err = g.GetNode(parent); err != nil {
		return err
	}

	if err = pip.CheckAssignment(childNode.Kind, parentNode.Kind); err != nil {
		return err
	}

	if _, ok := g.assignments[child]; !ok {
		g.assignments[child] = make(map[string]bool)
	}
	g.assignments[child][parent] = true

	return nil
}

func (g graph) Deassign(child string, parent string) error {
	delete(g.assignments[child], parent)
	return nil
}

func (g graph) GetChildren(name string) (map[string]pip.Node, error) {
	children := make(map[string]pip.Node)
	for nodeName, assignmentMap := range g.assignments {
		if assignmentMap[name] {
			node, _ := g.GetNode(nodeName)
			children[nodeName] = node
		}
	}
	return children, nil
}

func (g graph) GetParents(name string) (map[string]pip.Node, error) {
	assignments := g.assignments[name]
	parents := make(map[string]pip.Node)
	for nodeName := range assignments {
		node, _ := g.GetNode(nodeName)
		parents[nodeName] = node
	}
	return parents, nil
}

func (g graph) GetAssignments() (map[string]map[string]bool, error) {
	assignments := make(map[string]map[string]bool)
	for child, parents := range g.assignments {
		retParents := make(map[string]bool)
		for parent, ok := range parents {
			if !ok {
				continue
			}

			retParents[parent] = ok
		}

		assignments[child] = retParents
	}

	return assignments, nil
}

func (g graph) Associate(subject string, target string, operations pip.Operations) error {
	var (
		subjectNode pip.Node
		targetNode  pip.Node
		err         error
	)

	if subjectNode, err = g.GetNode(subject); err != nil {
		return err
	}
	if targetNode, err = g.GetNode(target); err != nil {
		return err
	}

	if err = pip.CheckAssociation(subjectNode.Kind, targetNode.Kind); err != nil {
		return err
	}

	if _, ok := g.associations[subject]; !ok {
		g.associations[subject] = make(map[string]pip.Operations)
	}
	g.associations[subject][target] = copyOps(operations)

	return nil
}

func (g graph) Dissociate(subject string, target string) error {
	delete(g.associations[subject], target)
	return nil
}

func (g graph) GetAssociationsForSubject(subject string) (map[string]pip.Operations, error) {
	retAssocs := make(map[string]pip.Operations)
	assocs := g.associations[subject]
	for target, ops := range assocs {
		retAssocs[target] = copyOps(ops)
	}
	return retAssocs, nil
}

func (g graph) GetAssociations() (map[string]map[string]pip.Operations, error) {
	assocs := make(map[string]map[string]pip.Operations)
	for subject, subjectAssocs := range g.associations {
		retAssocs := make(map[string]pip.Operations)
		for target, ops := range subjectAssocs {
			retAssocs[target] = copyOps(ops)
		}

		assocs[subject] = retAssocs
	}

	return assocs, nil
}

type jsonGraph struct {
	Nodes        map[string]pip.Node                  `json:"nodes"`
	Assignments  map[string]map[string]bool           `json:"assignments"`
	Associations map[string]map[string]pip.Operations `json:"associations"`
}

func (g graph) MarshalJSON() ([]byte, error) {
	var err error
	jg := jsonGraph{
		Nodes:        make(map[string]pip.Node),
		Assignments:  make(map[string]map[string]bool),
		Associations: make(map[string]map[string]pip.Operations),
	}

	if jg.Nodes, err = g.GetNodes(); err != nil {
		return nil, err
	}

	if jg.Assignments, err = g.GetAssignments(); err != nil {
		return nil, err
	}

	if jg.Associations, err = g.GetAssociations(); err != nil {
		return nil, err
	}

	return json.Marshal(jg)
}

// UnmarshalJSON into a graph.
// This will erase any nodes/assignments/associations that currently exist in the graph.
func (g *graph) UnmarshalJSON(bytes []byte) error {
	jg := jsonGraph{
		Nodes:        make(map[string]pip.Node),
		Assignments:  make(map[string]map[string]bool),
		Associations: make(map[string]map[string]pip.Operations),
	}

	if err := json.Unmarshal(bytes, &jg); err != nil {
		return err
	}

	g.nodes = jg.Nodes
	g.assignments = jg.Assignments
	g.associations = jg.Associations

	return nil
}

func copyNode(node pip.Node) pip.Node {
	props := make(map[string]string)
	for k, v := range node.Properties {
		props[k] = v
	}
	return pip.Node{
		Name:       node.Name,
		Kind:       node.Kind,
		Properties: props,
	}
}

func copyOps(operations pip.Operations) pip.Operations {
	retOps := pip.ToOps()
	for op := range operations {
		retOps[op] = true
	}
	return retOps
}

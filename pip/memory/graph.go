package memory

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/policy"
)

type (
	memgraph struct {
		nodes        map[string]policy.Node
		assignments  map[string]map[string]bool
		associations map[string]map[string]policy.Operations
	}
)

func NewGraph() policy.Graph {
	return &memgraph{
		nodes:        make(map[string]policy.Node),
		assignments:  make(map[string]map[string]bool),
		associations: make(map[string]map[string]policy.Operations),
	}
}

func (g *memgraph) CreatePolicyClass(name string) error {
	if _, ok := g.nodes[name]; ok {
		return fmt.Errorf("name %q already exists", name)
	}

	g.nodes[name] = policy.Node{
		Name:       name,
		Kind:       policy.PolicyClass,
		Properties: make(map[string]string),
	}

	return nil
}

func (g *memgraph) CreateNode(name string, kind policy.Kind, properties map[string]string, parent string, parents ...string) (policy.Node, error) {
	if _, ok := g.nodes[name]; ok {
		return policy.Node{}, fmt.Errorf("name %q already exists", name)
	}

	if properties == nil {
		properties = make(map[string]string)
	}

	// create the node
	n := policy.Node{
		Name:       name,
		Kind:       kind,
		Properties: properties,
	}
	node := copyNode(n)
	g.nodes[name] = node

	// set assignments for the new node
	assignments := make(map[string]bool)

	// check the initial parent exists
	if _, ok := g.nodes[parent]; !ok {
		return policy.Node{}, fmt.Errorf("parent %q does not exist", parent)
	}

	assignments[parent] = true

	// check other parents exist and add to assignments
	for _, p := range parents {
		if _, ok := g.nodes[parent]; !ok {
			return policy.Node{}, fmt.Errorf("parent %q does not exist", parent)
		}

		assignments[p] = true
	}

	g.assignments[name] = assignments

	return node, nil
}

func (g *memgraph) UpdateNode(name string, properties map[string]string) error {
	if ok, _ := g.Exists(name); !ok {
		return fmt.Errorf("node %q does not exist", name)
	}

	n := g.nodes[name]
	n.Properties = properties
	g.nodes[name] = copyNode(n)

	return nil
}

func (g *memgraph) DeleteNode(name string) error {
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

func (g *memgraph) Exists(name string) (bool, error) {
	_, ok := g.nodes[name]
	return ok, nil
}

func (g *memgraph) GetNodes() (map[string]policy.Node, error) {
	nodes := make(map[string]policy.Node)
	for _, node := range g.nodes {
		copyNode := copyNode(node)
		nodes[copyNode.Name] = copyNode
	}
	return nodes, nil
}

func (g *memgraph) GetNode(name string) (policy.Node, error) {
	node, ok := g.nodes[name]
	if !ok {
		return policy.Node{}, fmt.Errorf("node %q does not exist", name)
	}
	return copyNode(node), nil
}

func (g *memgraph) Find(kind policy.Kind, properties map[string]string) (map[string]policy.Node, error) {
	found := make(map[string]policy.Node)
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

func (g *memgraph) Assign(child string, parent string) error {
	var (
		childNode  policy.Node
		parentNode policy.Node
		err        error
	)

	if childNode, err = g.GetNode(child); err != nil {
		return err
	}
	if parentNode, err = g.GetNode(parent); err != nil {
		return err
	}

	if err = policy.CheckAssignment(childNode.Kind, parentNode.Kind); err != nil {
		return err
	}

	if _, ok := g.assignments[child]; !ok {
		g.assignments[child] = make(map[string]bool)
	}
	g.assignments[child][parent] = true

	return nil
}

func (g *memgraph) Deassign(child string, parent string) error {
	delete(g.assignments[child], parent)
	return nil
}

func (g *memgraph) GetChildren(name string) (map[string]policy.Node, error) {
	children := make(map[string]policy.Node)
	for nodeName, assignmentMap := range g.assignments {
		if assignmentMap[name] {
			node, _ := g.GetNode(nodeName)
			children[nodeName] = node
		}
	}
	return children, nil
}

func (g *memgraph) GetParents(name string) (map[string]policy.Node, error) {
	assignments := g.assignments[name]
	parents := make(map[string]policy.Node)
	for nodeName := range assignments {
		node, _ := g.GetNode(nodeName)
		parents[nodeName] = node
	}
	return parents, nil
}

func (g *memgraph) GetAssignments() (map[string]map[string]bool, error) {
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

func (g *memgraph) Associate(subject string, target string, operations policy.Operations) error {
	var (
		subjectNode policy.Node
		targetNode  policy.Node
		err         error
	)

	if subjectNode, err = g.GetNode(subject); err != nil {
		return err
	}
	if targetNode, err = g.GetNode(target); err != nil {
		return err
	}

	if err = policy.CheckAssociation(subjectNode.Kind, targetNode.Kind); err != nil {
		return err
	}

	if _, ok := g.associations[subject]; !ok {
		g.associations[subject] = make(map[string]policy.Operations)
	}
	g.associations[subject][target] = copyOps(operations)

	return nil
}

func (g *memgraph) Dissociate(subject string, target string) error {
	delete(g.associations[subject], target)
	return nil
}

func (g *memgraph) GetAssociationsForSubject(subject string) (map[string]policy.Operations, error) {
	retAssocs := make(map[string]policy.Operations)
	assocs := g.associations[subject]
	for target, ops := range assocs {
		retAssocs[target] = copyOps(ops)
	}
	return retAssocs, nil
}

func (g *memgraph) GetAssociations() (map[string]map[string]policy.Operations, error) {
	assocs := make(map[string]map[string]policy.Operations)
	for subject, subjectAssocs := range g.associations {
		retAssocs := make(map[string]policy.Operations)
		for target, ops := range subjectAssocs {
			retAssocs[target] = copyOps(ops)
		}

		assocs[subject] = retAssocs
	}

	return assocs, nil
}

type jsonGraph struct {
	Nodes        map[string]policy.Node                  `json:"nodes"`
	Assignments  map[string]map[string]bool              `json:"assignments"`
	Associations map[string]map[string]policy.Operations `json:"associations"`
}

func (g *memgraph) MarshalJSON() ([]byte, error) {
	var err error
	jg := jsonGraph{
		Nodes:        make(map[string]policy.Node),
		Assignments:  make(map[string]map[string]bool),
		Associations: make(map[string]map[string]policy.Operations),
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
func (g *memgraph) UnmarshalJSON(bytes []byte) error {
	jg := jsonGraph{
		Nodes:        make(map[string]policy.Node),
		Assignments:  make(map[string]map[string]bool),
		Associations: make(map[string]map[string]policy.Operations),
	}

	if err := json.Unmarshal(bytes, &jg); err != nil {
		return err
	}

	g.nodes = jg.Nodes
	g.assignments = jg.Assignments
	g.associations = jg.Associations

	return nil
}

func copyNode(node policy.Node) policy.Node {
	props := make(map[string]string)
	for k, v := range node.Properties {
		props[k] = v
	}
	return policy.Node{
		Name:       node.Name,
		Kind:       node.Kind,
		Properties: props,
	}
}

func copyOps(operations policy.Operations) policy.Operations {
	retOps := policy.ToOps()
	for op := range operations {
		retOps[op] = true
	}
	return retOps
}

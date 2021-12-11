package dag

import (
	"github.com/PM-Master/policy-machine-go/policy"
)

type (
	Propagator interface {
		Propagate(node policy.Node, target policy.Node)
	}

	Visitor interface {
		Visit(node policy.Node)
	}

	Searcher interface {
		Traverse(start policy.Node, propagate func(node policy.Node, target policy.Node) error, visit func(node policy.Node) error) error
	}

	bfs struct {
		graph policy.Graph
	}

	dfs struct {
		graph   policy.Graph
		visited map[string]bool
	}
)

func NewBFS(graph policy.Graph) Searcher {
	return bfs{graph: graph}
}

func NewDFS(graph policy.Graph) Searcher {
	return dfs{graph: graph, visited: make(map[string]bool)}
}

func (b bfs) Traverse(start policy.Node, propagate func(parent policy.Node, child policy.Node) error, visit func(node policy.Node) error) error {
	queue := make([]policy.Node, 0)
	queue = append(queue, start)

	seen := make(map[string]bool)
	seen[start.Name] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		var (
			err     error
			parents map[string]policy.Node
		)

		if err = visit(node); err != nil {
			return err
		}

		if parents, err = b.graph.GetParents(node.Name); err != nil {
			return err
		}

		for _, parentNode := range parents {
			if seen[parentNode.Name] {
				continue
			}

			queue = append(queue, parentNode)
			seen[parentNode.Name] = true

			if err = propagate(parentNode, node); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d dfs) Traverse(start policy.Node, propagate func(node policy.Node, target policy.Node) error, visit func(node policy.Node) error) error {
	if d.visited[start.Name] {
		return nil
	}

	d.visited[start.Name] = true

	parents, err := d.graph.GetParents(start.Name)
	if err != nil {
		return err
	}

	for _, parentNode := range parents {
		if err := d.Traverse(parentNode, propagate, visit); err != nil {
			return err
		}

		if err := propagate(parentNode, start); err != nil {
			return err
		}
	}

	if err := visit(start); err != nil {
		return nil
	}

	return nil
}

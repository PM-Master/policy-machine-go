package dag

import (
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
)

type (
	Propagator interface {
		Propagate(node graph.Node, target graph.Node)
	}

	Visitor interface {
		Visit(node graph.Node)
	}

	Searcher interface {
		Traverse(start graph.Node, propagate func(node graph.Node, target graph.Node) error, visit func(node graph.Node) error) error
	}

	bfs struct {
		graph ngac.Graph
	}

	dfs struct {
		graph   ngac.Graph
		visited map[string]bool
	}
)

func NewBFS(graph ngac.Graph) Searcher {
	return bfs{graph: graph}
}

func NewDFS(graph ngac.Graph) Searcher {
	return dfs{graph: graph, visited: make(map[string]bool)}
}

func (b bfs) Traverse(start graph.Node, propagate func(parent graph.Node, child graph.Node) error, visit func(node graph.Node) error) error {
	queue := make([]graph.Node, 0)
	queue = append(queue, start)

	seen := make(map[string]bool)
	seen[start.Name] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		var (
			err     error
			parents map[string]graph.Node
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

func (d dfs) Traverse(start graph.Node, propagate func(node graph.Node, target graph.Node) error, visit func(node graph.Node) error) error {
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

package dag

import (
	"github.com/PM-Master/policy-machine-go/pip"
)

type (
	Propagator interface {
		Propagate(node pip.Node, target pip.Node)
	}

	Visitor interface {
		Visit(node pip.Node)
	}

	Searcher interface {
		Traverse(start pip.Node, propagate func(node pip.Node, target pip.Node) error, visit func(node pip.Node) error) error
	}

	bfs struct {
		graph pip.Graph
	}

	dfs struct {
		graph   pip.Graph
		visited map[string]bool
	}
)

func NewBFS(graph pip.Graph) Searcher {
	return bfs{graph: graph}
}

func NewDFS(graph pip.Graph) Searcher {
	return dfs{graph: graph, visited: make(map[string]bool)}
}

func (b bfs) Traverse(start pip.Node, propagate func(parent pip.Node, child pip.Node) error, visit func(node pip.Node) error) error {
	queue := make([]pip.Node, 0)
	seen := make(map[string]bool)
	queue = append(queue, start)
	seen[start.Name] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		var (
			err     error
			parents map[string]pip.Node
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

func (d dfs) Traverse(start pip.Node, propagate func(node pip.Node, target pip.Node) error, visit func(node pip.Node) error) error {
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

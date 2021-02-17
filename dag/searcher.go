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
		Traverse(start pip.Node, propagate func(node pip.Node, target pip.Node), visit func(node pip.Node))
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

func (b bfs) Traverse(start pip.Node, propagate func(parent pip.Node, child pip.Node), visit func(node pip.Node)) {
	queue := make([]pip.Node, 0)
	seen := make(map[string]bool)
	queue = append(queue, start)
	seen[start.Name] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		visit(node)

		parents := b.graph.GetParents(node.Name)
		for _, parentNode := range parents {
			if seen[parentNode.Name] {
				continue
			}

			queue = append(queue, parentNode)
			seen[parentNode.Name] = true

			propagate(parentNode, node)
		}
	}
}

func (d dfs) Traverse(start pip.Node, propagate func(node pip.Node, target pip.Node), visit func(node pip.Node)) {
	if d.visited[start.Name] {
		return
	}

	d.visited[start.Name] = true

	parents := d.graph.GetParents(start.Name)
	for _, parentNode := range parents {
		d.Traverse(parentNode, propagate, visit)
		propagate(parentNode, start)
	}

	visit(start)
}

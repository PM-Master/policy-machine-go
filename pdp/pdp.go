package pdp

import (
	"github.com/PM-Master/policy-machine-go/dag"
	"github.com/PM-Master/policy-machine-go/pip"
)

type (
	Decider interface {
		Decide(user string, target string, permissions ...string) bool
	}

	decider struct {
		graph pip.Graph
	}

	userContext struct {
		borderTargets map[string]pip.Operations
	}

	targetContext struct {
		pcSet map[string]pip.Operations
	}
)

func NewDecider(graph pip.Graph) Decider {
	return decider{graph: graph}
}

func (d decider) Decide(user string, target string, permissions ...string) bool {
	// process user dag
	userNode, _ := d.graph.GetNode(user)
	userCtx := d.userDAG(userNode)

	// process target dag
	targetNode, _ := d.graph.GetNode(target)
	targetCtx := d.targetDAG(targetNode, userCtx)

	// resolve permissions
	allowed := d.allowedPermissions(targetCtx)
	for _, permission := range permissions {
		if !allowed[permission] {
			return false
		}
	}

	return true
}

func (d decider) userDAG(user pip.Node) userContext {
	bfs := dag.NewBFS(d.graph)
	userCtx := userContext{
		borderTargets: make(map[string]pip.Operations),
	}

	visitor := func(node pip.Node) {
		assocs := d.graph.GetAssociations(node.Name)
		d.collectAssociations(assocs, userCtx.borderTargets)
	}

	propagator := func(node pip.Node, parent pip.Node) {}

	bfs.Traverse(user, propagator, visitor)

	return userCtx
}

func (d decider) collectAssociations(assocs map[string]pip.Operations, borderTargets map[string]pip.Operations) {
	for target := range assocs {
		ops := assocs[target]
		if exOps, ok := borderTargets[target]; ok {
			for op := range exOps {
				ops[op] = true
			}
		}

		borderTargets[target] = ops
	}
}

func (d decider) targetDAG(target pip.Node, userCtx userContext) targetContext {
	visitedNodes := make(map[string]map[string]pip.Operations)

	visitor := func(node pip.Node) {
		nodeCtx, ok := visitedNodes[node.Name]
		if !ok {
			nodeCtx = make(map[string]pip.Operations)
			visitedNodes[node.Name] = nodeCtx
		}

		if node.Kind == pip.PolicyClass {
			nodeCtx[node.Name] = make(pip.Operations)
		} else {
			ops, ok := userCtx.borderTargets[node.Name]
			if ok {
				for pc := range nodeCtx {
					pcOps := nodeCtx[pc]
					for op := range ops {
						pcOps[op] = true
					}
					nodeCtx[pc] = pcOps
				}
			}
		}
	}

	propagator := func(parent pip.Node, child pip.Node) {
		parentCtx := visitedNodes[parent.Name]
		nodeCtx, ok := visitedNodes[child.Name]
		if !ok {
			nodeCtx = make(map[string]pip.Operations)
		}
		for name := range parentCtx {
			ops, ok := nodeCtx[name]
			if !ok {
				ops = make(pip.Operations)
			}

			parentOps := parentCtx[name]
			for op := range parentOps {
				ops[op] = true
			}

			nodeCtx[name] = ops
		}

		visitedNodes[child.Name] = nodeCtx
	}

	dfs := dag.NewDFS(d.graph)
	dfs.Traverse(target, propagator, visitor)

	return targetContext{pcSet: visitedNodes[target.Name]}
}

func (d decider) allowedPermissions(ctx targetContext) pip.Operations {
	allowed := make(pip.Operations)
	pcSet := ctx.pcSet
	first := true
	for _, ops := range pcSet {
		if first {
			for op := range ops {
				allowed[op] = true
			}
			first = false
		} else {
			if allowed[pip.AllOps] {
				allowed = make(pip.Operations)
				for op := range ops {
					allowed[op] = true
				}
			} else if !allowed[pip.AllOps] {
				for op := range allowed {
					if _, ok := ops[op]; !ok {
						delete(allowed, op)
					}
				}
			}
		}
	}

	return allowed
}

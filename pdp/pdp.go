package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/dag"
	"github.com/PM-Master/policy-machine-go/pip"
)

type (
	Decider interface {
		HasPermissions(user string, target string, permissions ...string) (bool, error)
		ListPermissions(user string, target string) (pip.Operations, error)
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

func (d decider) HasPermissions(user string, target string, permissions ...string) (bool, error) {
	allowed, err := d.ListPermissions(user, target)
	if err != nil {
		return false, fmt.Errorf("error checking if user %s has permissions %s on target %s", user, permissions, target)
	}

	for _, permission := range permissions {
		if !allowed.Contains(permission) {
			return false, nil
		}
	}

	return true, nil
}

func (d decider) ListPermissions(user string, target string) (pip.Operations, error) {
	var (
		userCtx   userContext
		targetCtx targetContext
		err       error
	)

	// process user dag
	userNode, _ := d.graph.GetNode(user)
	if userCtx, err = d.userDAG(userNode); err != nil {
		return nil, fmt.Errorf("error processing user side of graph for %q: %v", user, err)
	}

	// process target dag
	targetNode, _ := d.graph.GetNode(target)
	if targetCtx, err = d.targetDAG(targetNode, userCtx); err != nil {
		return nil, fmt.Errorf("error processing target side of graph for %q: %v", target, err)
	}

	// resolve permissions
	allowed := d.allowedPermissions(targetCtx)

	return allowed, nil
}

func (d decider) userDAG(user pip.Node) (userContext, error) {
	bfs := dag.NewBFS(d.graph)
	userCtx := userContext{
		borderTargets: make(map[string]pip.Operations),
	}

	visitor := func(node pip.Node) error {
		assocs, err := d.graph.GetAssociationsForSubject(node.Name)
		if err != nil {
			return err
		}
		d.collectAssociations(assocs, userCtx.borderTargets)

		return nil
	}

	propagator := func(node pip.Node, parent pip.Node) error {
		return nil
	}

	if err := bfs.Traverse(user, propagator, visitor); err != nil {
		return userContext{}, err
	}

	return userCtx, nil
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

func (d decider) targetDAG(target pip.Node, userCtx userContext) (targetContext, error) {
	visitedNodes := make(map[string]map[string]pip.Operations)

	visitor := func(node pip.Node) error {
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

		return nil
	}

	propagator := func(parent pip.Node, child pip.Node) error {
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
		return nil
	}

	dfs := dag.NewDFS(d.graph)
	err := dfs.Traverse(target, propagator, visitor)
	return targetContext{pcSet: visitedNodes[target.Name]}, err
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
			} else if !ops[pip.AllOps] {
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

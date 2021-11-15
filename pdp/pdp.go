package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/dag"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"github.com/PM-Master/policy-machine-go/pip/memory"
)

type (
	Decider interface {
		HasPermissions(user string, target string, permissions ...string) (bool, error)
		ListPermissions(user string, target string) (graph.Operations, error)
	}

	decider struct {
		graph        ngac.Graph
		prohibitions ngac.Prohibitions
	}

	userContext struct {
		borderTargets map[string]graph.Operations
		prohibitions  []ngac.Prohibition
	}

	targetContext struct {
		pcSet   map[string]graph.Operations
		visited map[string]bool
	}
)

func NewDecider(graph ngac.Graph, prohibitions ngac.Prohibitions) Decider {
	if prohibitions == nil {
		prohibitions = memory.NewProhibitions()
	}

	return decider{graph: graph, prohibitions: prohibitions}
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

func (d decider) ListPermissions(user string, target string) (graph.Operations, error) {
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
	allowed := d.resolvePermissions(userCtx, targetCtx, target)

	return allowed, nil
}

func (d decider) userDAG(user graph.Node) (userContext, error) {
	bfs := dag.NewBFS(d.graph)
	userCtx := userContext{
		borderTargets: make(map[string]graph.Operations),
		prohibitions:  make([]ngac.Prohibition, 0),
	}

	visitor := func(node graph.Node) error {
		assocs, err := d.graph.GetAssociationsForSubject(node.Name)
		if err != nil {
			return err
		}

		d.collectAssociations(assocs, userCtx.borderTargets)

		pros, err := d.prohibitions.Get(node.Name)
		if err != nil {
			return err
		}

		userCtx.prohibitions = append(userCtx.prohibitions, pros...)

		return nil
	}

	propagator := func(node graph.Node, parent graph.Node) error {
		return nil
	}

	if err := bfs.Traverse(user, propagator, visitor); err != nil {
		return userContext{}, err
	}

	return userCtx, nil
}

func (d decider) collectAssociations(assocs map[string]graph.Operations, borderTargets map[string]graph.Operations) {
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

func (d decider) targetDAG(target graph.Node, userCtx userContext) (targetContext, error) {
	foundPermissions := make(map[string]map[string]graph.Operations)
	visited := make(map[string]bool)

	visitor := func(node graph.Node) error {
		visited[node.Name] = true

		nodeCtx, ok := foundPermissions[node.Name]
		if !ok {
			nodeCtx = make(map[string]graph.Operations)
			foundPermissions[node.Name] = nodeCtx
		}

		if node.Kind == graph.PolicyClass {
			nodeCtx[node.Name] = make(graph.Operations)
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

	propagator := func(parent graph.Node, child graph.Node) error {
		parentCtx := foundPermissions[parent.Name]
		nodeCtx, ok := foundPermissions[child.Name]
		if !ok {
			nodeCtx = make(map[string]graph.Operations)
		}
		for name := range parentCtx {
			ops, ok := nodeCtx[name]
			if !ok {
				ops = make(graph.Operations)
			}

			parentOps := parentCtx[name]
			for op := range parentOps {
				ops[op] = true
			}

			nodeCtx[name] = ops
		}

		foundPermissions[child.Name] = nodeCtx
		return nil
	}

	dfs := dag.NewDFS(d.graph)
	err := dfs.Traverse(target, propagator, visitor)
	return targetContext{pcSet: foundPermissions[target.Name]}, err
}

func (d decider) resolvePermissions(userCtx userContext, targetContext targetContext, target string) graph.Operations {
	allowed := d.allowedPermissions(targetContext)
	denied := d.deniedPermissions(userCtx, targetContext, target)

	allowed.RemoveAll(denied)

	return allowed
}

func (d decider) allowedPermissions(ctx targetContext) graph.Operations {
	allowed := make(graph.Operations)
	pcSet := ctx.pcSet
	first := true
	for _, ops := range pcSet {
		if first {
			for op := range ops {
				allowed[op] = true
			}
			first = false
		} else {
			if allowed[graph.AllOps] {
				allowed = make(graph.Operations)
				for op := range ops {
					allowed[op] = true
				}
			} else if !ops[graph.AllOps] {
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

func (d decider) deniedPermissions(userCtx userContext, targetCtx targetContext, target string) graph.Operations {
	denied := make(graph.Operations)
	visited := targetCtx.visited
	prohibitions := userCtx.prohibitions

	for _, prohibition := range prohibitions {
		isIntersection := prohibition.Intersection
		containers := prohibition.Containers
		addOps := false

		for container, complement := range containers {
			if target == container {
				addOps = false
				if isIntersection {
					break
				} else {
					continue
				}
			}

			if (!complement && visited[target]) || (complement && !visited[target]) {
				addOps = true

				if !isIntersection {
					break
				}
			} else {
				addOps = false

				if isIntersection {
					break
				}
			}
		}

		if addOps {
			denied.AddAll(prohibition.Operations)
		}
	}

	return denied
}

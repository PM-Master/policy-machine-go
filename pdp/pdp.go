package pdp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/dag"
	"github.com/PM-Master/policy-machine-go/pip/memory"
	"github.com/PM-Master/policy-machine-go/policy"
)

type (
	Decider interface {
		HasPermissions(user string, target string, permissions ...string) (bool, error)
		ListPermissions(user string, target string) (policy.Operations, error)
	}

	decider struct {
		graph        policy.Graph
		prohibitions policy.Prohibitions
	}

	userContext struct {
		borderTargets map[string]policy.Operations
		prohibitions  []policy.Prohibition
	}

	targetContext struct {
		pcSet   map[string]policy.Operations
		visited map[string]bool
	}
)

func NewDecider(graph policy.Graph, prohibitions policy.Prohibitions) Decider {
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

func (d decider) ListPermissions(user string, target string) (policy.Operations, error) {
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

func (d decider) userDAG(user policy.Node) (userContext, error) {
	bfs := dag.NewBFS(d.graph)
	userCtx := userContext{
		borderTargets: make(map[string]policy.Operations),
		prohibitions:  make([]policy.Prohibition, 0),
	}

	visitor := func(node policy.Node) error {
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

	propagator := func(node policy.Node, parent policy.Node) error {
		return nil
	}

	if err := bfs.Traverse(user, propagator, visitor); err != nil {
		return userContext{}, err
	}

	return userCtx, nil
}

func (d decider) collectAssociations(assocs map[string]policy.Operations, borderTargets map[string]policy.Operations) {
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

func (d decider) targetDAG(target policy.Node, userCtx userContext) (targetContext, error) {
	foundPermissions := make(map[string]map[string]policy.Operations)
	visited := make(map[string]bool)

	visitor := func(node policy.Node) error {
		visited[node.Name] = true

		nodeCtx, ok := foundPermissions[node.Name]
		if !ok {
			nodeCtx = make(map[string]policy.Operations)
			foundPermissions[node.Name] = nodeCtx
		}

		if node.Kind == policy.PolicyClass {
			nodeCtx[node.Name] = make(policy.Operations)
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

	propagator := func(parent policy.Node, child policy.Node) error {
		parentCtx := foundPermissions[parent.Name]
		nodeCtx, ok := foundPermissions[child.Name]
		if !ok {
			nodeCtx = make(map[string]policy.Operations)
		}
		for name := range parentCtx {
			ops, ok := nodeCtx[name]
			if !ok {
				ops = make(policy.Operations)
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

func (d decider) resolvePermissions(userCtx userContext, targetContext targetContext, target string) policy.Operations {
	allowed := d.allowedPermissions(targetContext)
	denied := d.deniedPermissions(userCtx, targetContext, target)

	allowed.RemoveAll(denied)

	return allowed
}

func (d decider) allowedPermissions(ctx targetContext) policy.Operations {
	allowed := make(policy.Operations)
	pcSet := ctx.pcSet
	first := true
	for _, ops := range pcSet {
		if first {
			for op := range ops {
				allowed[op] = true
			}
			first = false
		} else {
			if allowed[policy.AllOps] {
				allowed = make(policy.Operations)
				for op := range ops {
					allowed[op] = true
				}
			} else if !ops[policy.AllOps] {
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

func (d decider) deniedPermissions(userCtx userContext, targetCtx targetContext, target string) policy.Operations {
	denied := make(policy.Operations)
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

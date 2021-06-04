package ngac

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
)

type (
	Graph interface {
		CreatePolicyClass(name string) error
		CreateNode(name string, kind graph.Kind, properties map[string]string, parent string, parents ...string) (graph.Node, error)
		UpdateNode(name string, properties map[string]string) error
		DeleteNode(name string) error
		Exists(name string) (bool, error)
		GetNodes() (map[string]graph.Node, error)
		GetNode(name string) (graph.Node, error)
		Find(kind graph.Kind, properties map[string]string) (map[string]graph.Node, error)
		Assign(child string, parent string) error
		Deassign(child string, parent string) error
		GetChildren(name string) (map[string]graph.Node, error)
		GetParents(name string) (map[string]graph.Node, error)
		GetAssignments() (map[string]map[string]bool, error)
		Associate(subject string, target string, operations graph.Operations) error
		Dissociate(subject string, target string) error
		GetAssociationsForSubject(subject string) (map[string]graph.Operations, error)
		GetAssociations() (map[string]map[string]graph.Operations, error)

		json.Marshaler
		json.Unmarshaler
	}
)

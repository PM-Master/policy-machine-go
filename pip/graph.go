package pip

import (
	"encoding/json"
)

type (
	Graph interface {
		CreateNode(name string, kind Kind, properties map[string]string) (Node, error)
		UpdateNode(name string, properties map[string]string) error
		DeleteNode(name string) error
		Exists(name string) (bool, error)
		GetNodes() (map[string]Node, error)
		GetNode(name string) (Node, error)
		Find(kind Kind, properties map[string]string) (map[string]Node, error)
		Assign(child string, parent string) error
		Deassign(child string, parent string) error
		GetChildren(name string) (map[string]Node, error)
		GetParents(name string) (map[string]Node, error)
		GetAssignments() (map[string]map[string]bool, error)
		Associate(subject string, target string, operations Operations) error
		Dissociate(subject string, target string) error
		GetAssociationsForSubject(subject string) (map[string]Operations, error)
		GetAssociations() (map[string]map[string]Operations, error)

		json.Marshaler
		json.Unmarshaler
	}
)

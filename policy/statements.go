package policy

import (
	"fmt"
	"strings"
)

type (
	Statement interface {
		Apply(store Store) error
	}

	CreatePolicyStatement struct {
		Name string `json:"name,omitempty"`
	}

	jsonCreatePolicyStatement struct {
		Name string `json:"name,omitempty"`
	}

	CreateNodeStatement struct {
		Name       string            `json:"name,omitempty"`
		Kind       Kind              `json:"kind,omitempty"`
		Properties map[string]string `json:"properties,omitempty"`
		Parents    []string          `json:"parents,omitempty"`
	}

	jsonCreateNodeStatement struct {
		Name       string            `json:"name,omitempty"`
		Kind       Kind              `json:"kind,omitempty"`
		Properties map[string]string `json:"properties,omitempty"`
		Parents    []string          `json:"parents,omitempty"`
	}

	AssignStatement struct {
		Child   string   `json:"child,omitempty"`
		Parents []string `json:"parents,omitempty"`
	}

	jsonAssignStatement struct {
		Child   string   `json:"child,omitempty"`
		Parents []string `json:"parents,omitempty"`
	}

	DeassignStatement struct {
		Child   string   `json:"child,omitempty"`
		Parents []string `json:"parents,omitempty"`
	}

	jsonDeassignStatement struct {
		Child   string   `json:"child,omitempty"`
		Parents []string `json:"parents,omitempty"`
	}

	DeleteNodeStatement struct {
		Name string `json:"name,omitempty"`
	}

	jsonDeleteNodeStatement struct {
		Name string `json:"name,omitempty"`
	}

	GrantStatement struct {
		Uattr      string     `json:"uattr,omitempty"`
		Target     string     `json:"target,omitempty"`
		Operations Operations `json:"operations,omitempty"`
	}

	jsonGrantStatement struct {
		Uattr      string     `json:"uattr,omitempty"`
		Target     string     `json:"target,omitempty"`
		Operations Operations `json:"operations,omitempty"`
	}

	DenyStatement struct {
		Subject      string     `json:"subject,omitempty"`
		Operations   Operations `json:"operations,omitempty"`
		Intersection bool       `json:"intersection,omitempty"`
		Containers   []string   `json:"containers,omitempty"`
	}

	jsonDenyStatement struct {
		Subject      string     `json:"subject,omitempty"`
		Operations   Operations `json:"operations,omitempty"`
		Intersection bool       `json:"intersection,omitempty"`
		Containers   []string   `json:"containers,omitempty"`
	}

	ObligationStatement struct {
		Obligation Obligation `json:"obligation"`
	}

	jsonObligationStatement struct {
		Obligation Obligation `json:"obligation"`
	}
)

func (c CreatePolicyStatement) Apply(store Store) error {
	err := store.Graph().CreatePolicyClass(c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c CreateNodeStatement) Apply(store Store) error {
	var err error

	if c.Kind == PolicyClass {
		err = store.Graph().CreatePolicyClass(c.Name)
	} else {
		_, err = store.Graph().CreateNode(c.Name, c.Kind, c.Properties, c.Parents[0], c.Parents[1:]...)
	}

	return err
}

func (d DeleteNodeStatement) Apply(store Store) error {
	return store.Graph().DeleteNode(d.Name)
}

func (a AssignStatement) Apply(store Store) error {
	for _, parent := range a.Parents {
		if err := store.Graph().Assign(a.Child, parent); err != nil {
			return fmt.Errorf("error assigning %s to %s", a.Child, parent)
		}
	}

	return nil
}

func (d DeassignStatement) Apply(store Store) error {
	for _, parent := range d.Parents {
		if err := store.Graph().Deassign(d.Child, parent); err != nil {
			return fmt.Errorf("error deassigning %s from %s", d.Child, parent)
		}
	}

	return nil
}

func (g GrantStatement) Apply(store Store) error {
	return store.Graph().Associate(g.Uattr, g.Target, g.Operations)
}

func (d DenyStatement) Apply(store Store) error {
	containers := make(map[string]bool)
	for _, containerName := range d.Containers {
		complement := strings.HasPrefix(containerName, "!")
		if complement {
			containerName = strings.TrimPrefix(containerName, "!")
		}

		containers[containerName] = complement
	}

	prohibition := Prohibition{
		Name:         fmt.Sprintf("deny-%s-%v-on-%v", d.Subject, d.Operations, d.Containers),
		Subject:      d.Subject,
		Containers:   containers,
		Operations:   d.Operations,
		Intersection: d.Intersection,
	}

	return store.Prohibitions().Add(prohibition)
}

func (o ObligationStatement) Apply(store Store) error {
	return store.Obligations().Add(o.Obligation)
}

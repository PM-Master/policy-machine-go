package ngac

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"strings"
)

type (
	Statement interface {
		Apply(fe FunctionalEntity) error
	}

	CreatePolicyStatement struct {
		Name string
	}

	CreateNodeStatement struct {
		Name       string
		Kind       graph.Kind
		Properties map[string]string
		Parents    []string
	}

	AssignStatement struct {
		Child   string
		Parents []string
	}

	DeassignStatement struct {
		Child   string
		Parents []string
	}

	DeleteNodeStatement struct {
		Name string
	}

	GrantStatement struct {
		Uattr      string
		Target     string
		Operations graph.Operations
	}

	DenyStatement struct {
		Subject      string
		Operations   graph.Operations
		Intersection bool
		Containers   []string
	}

	ObligationStatement struct {
		Obligation Obligation
	}
)

func (c CreatePolicyStatement) Apply(fe FunctionalEntity) error {
	err := fe.Graph().CreatePolicyClass(c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c CreateNodeStatement) Apply(fe FunctionalEntity) error {
	var err error

	if c.Kind == graph.PolicyClass {
		err = fe.Graph().CreatePolicyClass(c.Name)
	} else {
		_, err = fe.Graph().CreateNode(c.Name, c.Kind, c.Properties, c.Parents[0], c.Parents[1:]...)
	}

	return err
}

func (d DeleteNodeStatement) Apply(fe FunctionalEntity) error {
	return fe.Graph().DeleteNode(d.Name)
}

func (a AssignStatement) Apply(fe FunctionalEntity) error {
	for _, parent := range a.Parents {
		if err := fe.Graph().Assign(a.Child, parent); err != nil {
			return fmt.Errorf("error assigning %s to %s", a.Child, parent)
		}
	}

	return nil
}

func (d DeassignStatement) Apply(fe FunctionalEntity) error {
	for _, parent := range d.Parents {
		if err := fe.Graph().Deassign(d.Child, parent); err != nil {
			return fmt.Errorf("error deassigning %s from %s", d.Child, parent)
		}
	}

	return nil
}

func (g GrantStatement) Apply(fe FunctionalEntity) error {
	return fe.Graph().Associate(g.Uattr, g.Target, g.Operations)
}

func (d DenyStatement) Apply(fe FunctionalEntity) error {
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

	return fe.Prohibitions().Add(prohibition)
}

func (o ObligationStatement) Apply(fe FunctionalEntity) error {
	return fe.Obligations().Add(o.Obligation)
}

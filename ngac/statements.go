package ngac

import (
	"encoding/json"
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"strings"
)

type (
	Statement interface {
		Apply(fe FunctionalEntity) error

		json.Marshaler
		json.Unmarshaler
	}

	CreatePolicyStatement struct {
		Name string `json:"name,omitempty"`
	}

	jsonCreatePolicyStatement struct {
		Name string `json:"name,omitempty"`
	}

	CreateNodeStatement struct {
		Name       string            `json:"name,omitempty"`
		Kind       graph.Kind        `json:"kind,omitempty"`
		Properties map[string]string `json:"properties,omitempty"`
		Parents    []string          `json:"parents,omitempty"`
	}

	jsonCreateNodeStatement struct {
		Name       string            `json:"name,omitempty"`
		Kind       graph.Kind        `json:"kind,omitempty"`
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
		Uattr      string           `json:"uattr,omitempty"`
		Target     string           `json:"target,omitempty"`
		Operations graph.Operations `json:"operations,omitempty"`
	}

	jsonGrantStatement struct {
		Uattr      string           `json:"uattr,omitempty"`
		Target     string           `json:"target,omitempty"`
		Operations graph.Operations `json:"operations,omitempty"`
	}

	DenyStatement struct {
		Subject      string           `json:"subject,omitempty"`
		Operations   graph.Operations `json:"operations,omitempty"`
		Intersection bool             `json:"intersection,omitempty"`
		Containers   []string         `json:"containers,omitempty"`
	}

	jsonDenyStatement struct {
		Subject      string           `json:"subject,omitempty"`
		Operations   graph.Operations `json:"operations,omitempty"`
		Intersection bool             `json:"intersection,omitempty"`
		Containers   []string         `json:"containers,omitempty"`
	}

	ObligationStatement struct {
		Obligation Obligation `json:"obligation"`
	}

	jsonObligationStatement struct {
		Obligation Obligation `json:"obligation"`
	}
)

func (c *CreatePolicyStatement) Apply(fe FunctionalEntity) error {
	err := fe.Graph().CreatePolicyClass(c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c *CreatePolicyStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCreatePolicyStatement{
		Name: c.Name,
	})
}

func (c *CreatePolicyStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonCreatePolicyStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	c.Name = j.Name

	return nil
}

func (c *CreateNodeStatement) Apply(fe FunctionalEntity) error {
	var err error

	if c.Kind == graph.PolicyClass {
		err = fe.Graph().CreatePolicyClass(c.Name)
	} else {
		_, err = fe.Graph().CreateNode(c.Name, c.Kind, c.Properties, c.Parents[0], c.Parents[1:]...)
	}

	return err
}

func (c *CreateNodeStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCreateNodeStatement{
		Name:       c.Name,
		Kind:       c.Kind,
		Properties: c.Properties,
		Parents:    c.Parents,
	})
}

func (c *CreateNodeStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonCreateNodeStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	c.Name = j.Name
	c.Kind = j.Kind
	c.Properties = j.Properties
	c.Parents = j.Parents

	return nil
}

func (d *DeleteNodeStatement) Apply(fe FunctionalEntity) error {
	return fe.Graph().DeleteNode(d.Name)
}

func (d *DeleteNodeStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonDeleteNodeStatement{Name: d.Name})
}

func (d *DeleteNodeStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonDeleteNodeStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	d.Name = j.Name

	return nil
}

func (a *AssignStatement) Apply(fe FunctionalEntity) error {
	for _, parent := range a.Parents {
		if err := fe.Graph().Assign(a.Child, parent); err != nil {
			return fmt.Errorf("error assigning %s to %s", a.Child, parent)
		}
	}

	return nil
}

func (a *AssignStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAssignStatement{
		Child:   a.Child,
		Parents: a.Parents,
	})
}

func (a *AssignStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonAssignStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	a.Child = j.Child
	a.Parents = j.Parents

	return nil
}

func (d *DeassignStatement) Apply(fe FunctionalEntity) error {
	for _, parent := range d.Parents {
		if err := fe.Graph().Deassign(d.Child, parent); err != nil {
			return fmt.Errorf("error deassigning %s from %s", d.Child, parent)
		}
	}

	return nil
}

func (d *DeassignStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonDeassignStatement{
		Child:   d.Child,
		Parents: d.Parents,
	})
}

func (d *DeassignStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonDeassignStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	d.Child = j.Child
	d.Parents = j.Parents

	return nil
}

func (g *GrantStatement) Apply(fe FunctionalEntity) error {
	return fe.Graph().Associate(g.Uattr, g.Target, g.Operations)
}

func (g *GrantStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonGrantStatement{
		Uattr:      g.Uattr,
		Target:     g.Target,
		Operations: g.Operations,
	})
}

func (g *GrantStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonGrantStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	g.Uattr = j.Uattr
	g.Target = j.Target
	g.Operations = j.Operations

	return nil
}

func (d *DenyStatement) Apply(fe FunctionalEntity) error {
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

func (d *DenyStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonDenyStatement{
		Subject:      d.Subject,
		Operations:   d.Operations,
		Intersection: d.Intersection,
		Containers:   d.Containers,
	})
}

func (d *DenyStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonDenyStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	d.Subject = j.Subject
	d.Operations = j.Operations
	d.Intersection = j.Intersection
	d.Containers = j.Containers

	return nil
}

func (o *ObligationStatement) Apply(fe FunctionalEntity) error {
	return fe.Obligations().Add(o.Obligation)
}

func (o *ObligationStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonObligationStatement{
		Obligation: o.Obligation,
	})
}

func (o *ObligationStatement) UnmarshalJSON(bytes []byte) error {
	j := &jsonObligationStatement{}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}

	o.Obligation = j.Obligation

	return nil
}

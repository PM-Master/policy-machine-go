package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/policy"
	"io/ioutil"
	"os"
)

type Properties map[string]string

func ReadAndApply(policyStore policy.Store, path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("recevied directory path, expected .ngac file")
	}

	var pal []byte
	if pal, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("error reading file %q: %w", fileInfo.Name(), err)
	}

	return apply(policyStore, string(pal))
}

func apply(policyStore policy.Store, pal string) error {
	stmts, _, err := Parse(pal)
	if err != nil {
		return fmt.Errorf("error parsing policy author language: %w", err)
	}

	for _, stmt := range stmts {
		err = stmt.Apply(policyStore)
		if err != nil {
			return fmt.Errorf("error applying statement: %w", err)
		}
	}

	return nil
}

func Author(policyStore policy.Store, stmts ...policy.Statement) error {
	for _, stmt := range stmts {
		if err := stmt.Apply(policyStore); err != nil {
			return err
		}
	}

	return nil
}

func PolicyClass(name string) *policy.CreatePolicyStatement {
	return &policy.CreatePolicyStatement{
		Name: name,
	}
}

func ObjectAttribute(name string, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.ObjectAttribute,
		Properties: Properties{},
		Parents:    parents,
	}
}

func ObjectAttributeWithProperties(name string, properties Properties, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.ObjectAttribute,
		Properties: properties,
		Parents:    parents,
	}
}

func UserAttribute(name string, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.UserAttribute,
		Properties: Properties{},
		Parents:    parents,
	}
}

func UserAttributeWithProperties(name string, properties Properties, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.UserAttribute,
		Properties: properties,
		Parents:    parents,
	}
}

func Object(name string, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.Object,
		Properties: Properties{},
		Parents:    parents,
	}
}

func ObjectWithProperties(name string, properties Properties, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.Object,
		Properties: properties,
		Parents:    parents,
	}
}

func User(name string, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.User,
		Properties: Properties{},
		Parents:    parents,
	}
}

func UserWithProperties(name string, properties Properties, parent string, parents ...string) *policy.CreateNodeStatement {
	parents = append(parents, parent)
	return &policy.CreateNodeStatement{
		Name:       name,
		Kind:       policy.User,
		Properties: properties,
		Parents:    parents,
	}
}

func Assign(child string, parents ...string) *policy.AssignStatement {
	return &policy.AssignStatement{
		Child:   child,
		Parents: parents,
	}
}

func Deassign(child string, parents ...string) *policy.DeassignStatement {
	return &policy.DeassignStatement{
		Child:   child,
		Parents: parents,
	}
}

func Grant(ua string, target string, ops ...string) *policy.GrantStatement {
	return &policy.GrantStatement{
		Uattr:      ua,
		Target:     target,
		Operations: policy.ToOps(ops...),
	}
}

func Deny(subject string, operations policy.Operations, intersection bool, containers ...string) *policy.DenyStatement {
	return &policy.DenyStatement{
		Subject:      subject,
		Operations:   operations,
		Intersection: intersection,
		Containers:   containers,
	}
}

func Obligation(label string, event policy.EventPattern, response policy.ResponsePattern) *policy.ObligationStatement {
	return &policy.ObligationStatement{
		Obligation: policy.Obligation{
			User:     "",
			Label:    label,
			Event:    event,
			Response: response,
		},
	}
}

func Event(subject policy.Subject, operations []policy.EventOperation, containers ...string) policy.EventPattern {
	return policy.EventPattern{
		Subject:    subject,
		Operations: operations,
		Containers: containers,
	}
}

func Response(actions ...policy.Statement) policy.ResponsePattern {
	return policy.ResponsePattern{
		Actions: actions,
	}
}

func Operation(op string, args ...string) []policy.EventOperation {
	return []policy.EventOperation{
		{op, args},
	}
}

func Operations(ops ...string) []policy.EventOperation {
	eventOps := make([]policy.EventOperation, 0)
	for _, op := range ops {
		eventOps = append(eventOps, policy.EventOperation{
			Operation: op,
			Args:      []string{},
		})
	}

	return eventOps
}

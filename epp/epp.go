package epp

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"strings"
)

type (
	EventProcessor interface {
		ProcessEvent(eventCtx EventContext) error
	}

	EventContext struct {
		User   string
		Event  string
		Target string
		Args   map[string]string
	}

	epp struct {
		pap ngac.FunctionalEntity
	}
)

func NewEPP(pap ngac.FunctionalEntity) EventProcessor {
	return epp{pap: pap}
}

func (e epp) ProcessEvent(eventCtx EventContext) error {
	obligations, err := e.pap.Obligations().All()
	if err != nil {
		return fmt.Errorf("error getting obligations from PAP")
	}

	for _, obligation := range obligations {
		var matches bool
		matches, err = eventCtx.Matches(obligation.Event)
		if err != nil {
			return fmt.Errorf("error matching event pattern")
		}

		if !matches {
			continue
		}

		for _, action := range obligation.Response.Actions {
			action, err = resolveArgs(action, eventCtx.Args)
			if err != nil {
				return fmt.Errorf("error resolving args: %w", err)
			}

			err = action.Apply(e.pap)
			if err != nil {
				return fmt.Errorf("error applying response action: %w", err)
			}
		}
	}

	return nil
}

func resolveArgs(stmt ngac.Statement, args map[string]string) (ngac.Statement, error) {
	var err error

	if createPCStmt, ok := stmt.(*ngac.CreatePolicyStatement); ok {
		createPCStmt.Name = replaceArgs(createPCStmt.Name, args)
		return createPCStmt, nil
	} else if createNodeStmt, ok := stmt.(*ngac.CreateNodeStatement); ok {
		createNodeStmt.Name = replaceArgs(createNodeStmt.Name, args)

		// resolve properties
		for k, v := range createNodeStmt.Properties {
			createNodeStmt.Properties[k] = replaceArgs(v, args)
		}

		// resolve parents
		createNodeStmt.Parents = resolveSlice(createNodeStmt.Parents, args)

		return createNodeStmt, nil
	} else if assignStmt, ok := stmt.(*ngac.AssignStatement); ok {
		assignStmt.Child = replaceArgs(assignStmt.Child, args)
		assignStmt.Parents = resolveSlice(assignStmt.Parents, args)

		return assignStmt, nil
	} else if deassignStmt, ok := stmt.(*ngac.DeassignStatement); ok {
		deassignStmt.Child = replaceArgs(deassignStmt.Child, args)
		deassignStmt.Parents = resolveSlice(deassignStmt.Parents, args)

		return deassignStmt, nil
	} else if deleteNodeStmt, ok := stmt.(*ngac.DeleteNodeStatement); ok {
		deleteNodeStmt.Name = replaceArgs(deleteNodeStmt.Name, args)

		return deleteNodeStmt, nil
	} else if grantStmt, ok := stmt.(*ngac.GrantStatement); ok {
		grantStmt.Uattr = replaceArgs(grantStmt.Uattr, args)
		grantStmt.Target = replaceArgs(grantStmt.Target, args)

		return grantStmt, nil
	} else if denyStmt, ok := stmt.(*ngac.DenyStatement); ok {
		denyStmt.Subject = replaceArgs(denyStmt.Subject, args)
		denyStmt.Containers = resolveSlice(denyStmt.Containers, args)

		return denyStmt, nil
	} else if oblStmt, ok := stmt.(*ngac.ObligationStatement); ok {
		oblStmt.Obligation.Label = replaceArgs(oblStmt.Obligation.Label, args)
		oblStmt.Obligation.Response.Actions, err = resolveStatements(oblStmt.Obligation.Response.Actions, args)
		if err != nil {
			return nil, fmt.Errorf("error resolving args for obligation statement: %w", err)
		}

		return oblStmt, nil
	} else {
		return nil, fmt.Errorf("unknown statement: %v", stmt)
	}
}

func resolveStatements(statements []ngac.Statement, args map[string]string) ([]ngac.Statement, error) {
	resolved := make([]ngac.Statement, 0)
	for _, s := range statements {
		statement, err := resolveArgs(s, args)
		if err != nil {
			return nil, fmt.Errorf("error resolving args: %w", err)
		}
		resolved = append(resolved, statement)
	}
	return resolved, nil
}

func resolveSlice(slice []string, args map[string]string) []string {
	resolved := make([]string, 0)
	for _, cont := range slice {
		cont = replaceArgs(cont, args)
		resolved = append(resolved, cont)
	}

	return resolved
}

func replaceArgs(str string, args map[string]string) string {
	for arg, value := range args {
		if strings.Contains(str, arg) {
			str = strings.ReplaceAll(str, fmt.Sprintf("$%s", arg), value)
		}
	}

	return str
}

func (e EventContext) Matches(eventPattern ngac.EventPattern) (bool, error) {
	var eventMatches bool
	for _, patternEvent := range eventPattern.Operations {
		if e.Event == patternEvent.Operation {
			eventMatches = true
			break
		}
	}

	if !eventMatches {
		return false, nil
	}

	return e.subjectMatches(e.User, eventPattern.Subject) && e.targetMatches(e.Target, eventPattern.Containers), nil
}

func (e EventContext) subjectMatches(eventUser string, patternSubject string) bool {
	if patternSubject == "ANY_USER" {
		return true
	}

	return strings.ToUpper(patternSubject) == "ANY_USER" || eventUser == patternSubject
}

func (e EventContext) targetMatches(eventTarget string, patternContainers []string) bool {
	if len(patternContainers) == 0 {
		return true
	}

	for _, cont := range patternContainers {
		if cont == eventTarget {
			return true
		}
	}

	return false
}

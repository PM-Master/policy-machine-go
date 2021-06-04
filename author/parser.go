package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"strings"
)

func Parse(pal string) ([]ngac.Statement, error) {
	statements := splitStatements(pal)
	return parseStatements(statements)
}

func parseStatements(statements []string) ([]ngac.Statement, error) {
	stmts := make([]ngac.Statement, 0)
	for _, stmtStr := range statements {
		stmtStr = strings.TrimSpace(stmtStr)
		stmtStr = strings.TrimSuffix(stmtStr, ";")
		upperStmtStr := strings.ToUpper(stmtStr)

		var (
			stmt ngac.Statement
			err  error
		)

		if strings.HasPrefix(upperStmtStr, "CREATE POLICY") {
			stmt, err = parseCreatePolicy(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "CREATE") {
			stmt, err = parseCreateNode(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "OBLIGATION") {
			obligationParser := NewObligationParser()
			o, err := obligationParser.Parse(stmtStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing obligation: %w", err)
			}

			stmt = ngac.ObligationStatement{Obligation: o}
		} else if strings.HasPrefix(upperStmtStr, "ASSIGN") {
			stmt, err = parseAssign(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "DEASSIGN") {
			stmt, err = parseDeassign(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "DELETE") {
			stmt, err = parseDelete(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "GRANT") {
			stmt, err = parseGrant(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "DENY") {
			stmt, err = parseDeny(stmtStr)
		}

		if err != nil {
			return nil, fmt.Errorf("error parsing statement %q: %w", stmtStr, err)
		}

		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func parseDeny(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	subject := fields[1]

	index := 2
	field := ""
	for index, field = range fields {
		if strings.ToUpper(field) == "ON" ||
			strings.ToUpper(field) == "INTERSECTION" {
			break
		}
	}

	permStr := strings.Join(fields[2:index], " ")
	split := strings.Split(permStr, ",")
	ops := make(graph.Operations, 0)
	for _, s := range split {
		ops.Add(strings.TrimSpace(s))
	}

	stmtStr = strings.Join(fields[index:], " ")

	if !strings.HasPrefix(strings.ToUpper(stmtStr), "ON") {
		return nil, fmt.Errorf("DENY statement must have an ON clause")
	}

	inter := false
	if strings.HasPrefix(strings.ToUpper(stmtStr), "ON INTERSECTION OF") {
		inter = true
		stmtStr = strings.Join(fields[(index+3):], " ")
	} else {
		stmtStr = strings.Join(fields[(index+1):], " ")
	}

	split = strings.Split(stmtStr, ",")
	containers := make([]string, 0)
	for _, s := range split {
		s = strings.TrimSpace(s)
		containers = append(containers, s)
	}

	return ngac.DenyStatement{
		Subject:      subject,
		Operations:   ops,
		Intersection: inter,
		Containers:   containers,
	}, nil
}

//`GRANT <user_Attribute> {<permission>} ON <user_or_object_attribute>;`
func parseGrant(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	uattr := fields[1]

	index := 2
	field := ""
	for index, field = range fields {
		if strings.ToUpper(field) == "ON" {
			break
		}
	}

	permStr := strings.Join(fields[2:index], " ")
	split := strings.Split(permStr, ",")
	ops := make(graph.Operations, 0)
	for _, s := range split {
		ops.Add(strings.TrimSpace(s))
	}

	stmtStr = strings.Join(fields[index:], " ")

	if !strings.HasPrefix(strings.ToUpper(stmtStr), "ON") {
		return nil, fmt.Errorf("GRANT statement must have an ON clause")
	}

	fields = strings.Fields(stmtStr)
	target := fields[1]

	return ngac.GrantStatement{
		Uattr:      uattr,
		Target:     target,
		Operations: ops,
	}, nil
}

func parseDelete(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	target := fields[2]
	return ngac.DeleteNodeStatement{
		Name: target,
	}, nil
}

// `DEASSIGN <child> FROM {<parent>};`
func parseDeassign(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	child := fields[1]

	parentsStr := strings.Join(fields[3:], " ")
	split := strings.Split(parentsStr, ",")
	parents := make([]string, 0)

	for _, s := range split {
		s = strings.TrimSpace(s)
		parents = append(parents, s)
	}

	return ngac.DeassignStatement{
		Child:   child,
		Parents: parents,
	}, nil
}

// `ASSIGN <child> TO {<parent>};`
func parseAssign(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	child := fields[1]

	parentsStr := strings.Join(fields[3:], " ")
	split := strings.Split(parentsStr, ",")
	parents := make([]string, 0)

	for _, s := range split {
		s = strings.TrimSpace(s)
		parents = append(parents, s)
	}

	return ngac.AssignStatement{
		Child:   child,
		Parents: parents,
	}, nil
}

func parseCreateNode(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)

	kindField := fields[1]
	attrOrNameField := fields[2]
	var (
		name     string
		endIndex int
	)

	if strings.EqualFold(attrOrNameField, "attribute") {
		kindField = fmt.Sprintf("%v %v", kindField, attrOrNameField)
		name = fields[3]
		endIndex = 4
	} else {
		name = attrOrNameField
		endIndex = 3
	}

	var kind graph.Kind
	switch strings.ToUpper(kindField) {
	case "USER ATTRIBUTE":
		kind = graph.UserAttribute
	case "OBJECT ATTRIBUTE":
		kind = graph.ObjectAttribute
	case "OBJECT":
		kind = graph.Object
	case "USER":
		kind = graph.User
	}

	properties := make(map[string]string)
	stmtStr = strings.Join(fields[endIndex:], " ")
	if strings.Contains(strings.ToUpper(stmtStr), "WITH PROPERTIES") {
		propFields := strings.Fields(stmtStr)
		endIndex = 0
		var f string
		for endIndex, f = range propFields {
			if f == "IN" {
				break
			}
		}
		endIndex--

		propsStr := strings.Join(propFields[2:endIndex], " ")
		split := strings.Split(propsStr, ",")
		for _, prop := range split {
			prop = strings.TrimSpace(prop)
			kv := strings.Split(prop, "=")
			properties[kv[0]] = kv[1]
		}

		stmtStr = strings.Join(propFields[endIndex:], " ")
	}

	if !strings.HasPrefix(strings.ToUpper(stmtStr), "IN") {
		return nil, fmt.Errorf("IN clause required for creating nodes")
	}

	// remove IN
	stmtStr = strings.TrimSuffix(stmtStr[3:], ";")

	// split parents by comma
	parents := make([]string, 0)
	split := strings.Split(stmtStr, ",")
	for _, s := range split {
		s = strings.TrimSpace(s)
		parents = append(parents, s)
	}

	return ngac.CreateNodeStatement{
		Name:       name,
		Kind:       kind,
		Properties: properties,
		Parents:    parents,
	}, nil
}

func parseCreatePolicy(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	name := strings.ReplaceAll(fields[2], "(", "")
	startIndex := strings.Index(stmtStr, "(") + 1
	endIndex := strings.LastIndex(stmtStr, ")")
	stmtsStr := strings.TrimSpace(stmtStr[startIndex:endIndex])

	stmts, err := Parse(stmtsStr)
	if err != nil {
		return nil, err
	}

	return ngac.CreatePolicyStatement{
		Name:       name,
		Statements: stmts,
	}, nil
}

func splitStatements(pal string) []string {
	stmts := make([]string, 0)
	stmt := ""
	parenCounter := 0
	fields := strings.Fields(pal)
	for _, f := range fields {
		stmt = fmt.Sprintf("%v %v", stmt, f)
		if strings.Contains(f, ";") {
			if strings.Contains(f, ")") {
				parenCounter--
			}

			if parenCounter != 0 {
				continue
			}

			stmts = append(stmts, stmt)
			stmt = ""
		} else if strings.Contains(f, "(") {
			parenCounter++
		} else if strings.Contains(f, ")") {
			parenCounter--
		}
	}

	return stmts
}

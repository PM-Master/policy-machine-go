package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
	"strings"
)

type ParsedFunction struct {
	Name  string
	Args  map[string]bool
	Stmts string
}

func Parse(pal string) ([]ngac.Statement, map[string]ParsedFunction, error) {
	split := strings.Split(pal, "\n")
	lines := make([]string, 0)
	for _, s := range split {
		if strings.HasPrefix(strings.TrimSpace(s), "#") {
			continue
		}

		lines = append(lines, s)
	}

	pal = strings.Join(lines, "\n")

	statements := splitStatements(pal)
	return parseStatements(statements)
}

func parseStatements(statements []string) ([]ngac.Statement, map[string]ParsedFunction, error) {
	stmts := make([]ngac.Statement, 0)
	functions := make(map[string]ParsedFunction)
	vars := make(map[string]string, 0)
	for _, stmtStr := range statements {
		stmtStr = strings.TrimSuffix(strings.TrimSpace(stmtStr), ";")
		upperStmtStr := strings.ToUpper(stmtStr)

		// replace variables in stmtStr
		stmtStr = resolveVars(stmtStr, vars)

		var (
			stmt ngac.Statement
			err  error
		)

		/*if strings.HasPrefix(upperStmtStr, "CREATE POLICY") {
			stmt, err = parseCreatePolicy(stmtStr)
		} else */
		if strings.HasPrefix(upperStmtStr, "CREATE") {
			stmt, err = parseCreateNode(stmtStr)
		} else if strings.HasPrefix(upperStmtStr, "OBLIGATION") {
			obligationParser := NewObligationParser()
			o, err := obligationParser.Parse(stmtStr)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing obligation: %w", err)
			}

			stmt = &ngac.ObligationStatement{Obligation: o}
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
		} else if strings.HasPrefix(upperStmtStr, "FUNC") {
			function, err := parseFunc(stmtStr)
			if err != nil {
				return nil, nil, err
			}

			functions[function.Name] = function

			continue
		} else if strings.HasPrefix(upperStmtStr, "#") {
			continue
		} else if strings.HasPrefix(upperStmtStr, "LET") {
			varName, varValue, err := parseVar(stmtStr)
			if err != nil {
				return nil, nil, err
			}

			varName = fmt.Sprintf("$%s", varName)
			vars[varName] = varValue

			continue
		} else {
			err = fmt.Errorf("unknown statement %s", stmtStr)
		}

		if err != nil {
			return nil, nil, fmt.Errorf("error parsing statement %q: %w", stmtStr, err)
		}

		stmts = append(stmts, stmt)
	}

	return stmts, functions, nil
}

func resolveVars(stmt string, vars map[string]string) string {
	for name, value := range vars {
		stmt = strings.ReplaceAll(stmt, name, value)
	}

	return stmt
}

func parseVar(varStr string) (string, string, error) {
	fields := strings.Fields(varStr)
	if len(fields) < 4 {
		return "", "", fmt.Errorf("not enough tokens in variable decleration statement %q", varStr)
	}

	return fields[1], fields[3], nil
}

func parseFunc(funcStr string) (ParsedFunction, error) {
	funcStr = strings.TrimSpace(strings.Replace(funcStr, "func", "", 1))
	index := strings.Index(funcStr, "{")
	funcDefStr := funcStr[0:index]
	index = strings.Index(funcStr, "(")
	funcName := funcDefStr[0:index]

	argsStr := funcStr[strings.Index(funcStr, "(")+1 : strings.Index(funcStr, ")")]
	argFields := strings.Fields(argsStr)
	args := make(map[string]bool, 0)
	for _, argField := range argFields {
		argField = strings.TrimSpace(strings.TrimSuffix(argField, ","))
		args[argField] = true
	}

	stmtStr := funcStr[strings.Index(funcStr, "{")+1 : strings.Index(funcStr, "}")]

	return ParsedFunction{
		Name:  funcName,
		Args:  args,
		Stmts: stmtStr,
	}, nil
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

	return &ngac.DenyStatement{
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

	return &ngac.GrantStatement{
		Uattr:      uattr,
		Target:     target,
		Operations: ops,
	}, nil
}

func parseDelete(stmtStr string) (ngac.Statement, error) {
	fields := strings.Fields(stmtStr)
	target := fields[1]
	return &ngac.DeleteNodeStatement{
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

	return &ngac.DeassignStatement{
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

	return &ngac.AssignStatement{
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
	case "POLICY":
		kind = graph.PolicyClass
		return &ngac.CreateNodeStatement{
			Name:       name,
			Kind:       kind,
			Properties: make(map[string]string),
			Parents:    make([]string, 0),
		}, nil
	case "USER ATTRIBUTE":
		kind = graph.UserAttribute
	case "OBJECT ATTRIBUTE":
		kind = graph.ObjectAttribute
	case "OBJECT":
		kind = graph.Object
	case "USER":
		kind = graph.User
	}

	stmtStr = strings.Join(fields[endIndex:], " ")
	properties := make(map[string]string)
	if strings.Contains(strings.ToUpper(stmtStr), "WITH PROPERTIES") {
		propFields := strings.Fields(stmtStr)
		endIndex = 0
		var f string
		for endIndex, f = range propFields {
			if strings.ToUpper(f) == "IN" {
				break
			}
		}

		propsStr := strings.Join(propFields[2:endIndex], " ")
		split := strings.Split(propsStr, ",")
		for _, prop := range split {
			if !strings.Contains(prop, "=") {
				continue
			}

			prop = strings.TrimSpace(prop)
			kv := strings.Split(prop, "=")
			properties[kv[0]] = kv[1]
		}

		stmtStr = strings.Join(propFields[endIndex:], " ")
	}

	if !strings.HasPrefix(strings.ToUpper(stmtStr), "IN") {
		return nil, fmt.Errorf("IN clause required for creating non policy class nodes")
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

	return &ngac.CreateNodeStatement{
		Name:       name,
		Kind:       kind,
		Properties: properties,
		Parents:    parents,
	}, nil
}

func splitStatements(pal string) []string {
	stmts := make([]string, 0)
	split := strings.Split(pal, ";")
	parenCounter := 0
	stmt := ""
	for _, s := range split {
		// add the semi colon back which will help with obligation sub statements
		s = strings.TrimSpace(s) + ";"
		if len(s) == 1 {
			continue
		}

		stmt += s

		for _, c := range s {
			if c == '(' {
				parenCounter++
			} else if c == ')' {
				parenCounter--
			}
		}

		if parenCounter == 0 {
			stmts = append(stmts, stmt)
			stmt = ""
		}
	}

	return stmts
}

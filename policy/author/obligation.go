package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/policy"
	"strings"
)

type (
	ObligationParser interface {
		Parse(obligation string) (policy.Obligation, error)
	}

	EventParser interface {
		Parse(event string) (policy.EventPattern, error)
	}

	ResponseParser interface {
		Parse(response string) (policy.ResponsePattern, error)
	}

	obligationParser struct {
		eventParser    EventParser
		responseParser ResponseParser
	}

	eventParser struct {
	}

	responseParser struct {
	}
)

const (
	ObligationKeyword = "OBLIGATION"
	When              = "WHEN"
	Performs          = "PERFORMS"
	On                = "ON"
	Do                = "DO"
	Or                = "OR"
)

func NewObligationParser() ObligationParser {
	return obligationParser{
		eventParser:    eventParser{},
		responseParser: responseParser{},
	}
}

func (o obligationParser) Parse(obligation string) (policy.Obligation, error) {
	fields := strings.Fields(obligation)

	label := fields[1]
	index := 2
	for index = range fields {
		if strings.HasPrefix(strings.ToUpper(fields[index]), Do) {
			break
		}
	}

	event := strings.Join(fields[2:index], " ")
	eventPattern, err := o.eventParser.Parse(event)
	if err != nil {
		return policy.Obligation{}, fmt.Errorf("error parsing event: %w", err)
	}

	response := strings.Join(fields[index:], " ")
	responsePattern, err := o.responseParser.Parse(response)
	if err != nil {
		return policy.Obligation{}, fmt.Errorf("error parsing response: %w", err)
	}

	return policy.Obligation{
		Label:    label,
		Event:    eventPattern,
		Response: responsePattern,
	}, nil
}

func (e eventParser) Parse(event string) (policy.EventPattern, error) {
	fields := strings.Fields(event)

	subject := fields[1]

	var index int
	for index = 3; index < len(fields); index++ {
		if strings.HasPrefix(strings.ToUpper(fields[index]), On) {
			break
		}
	}

	performs := fields[3:index]
	ops, err := e.parsePerforms(strings.Join(performs, " "))
	if err != nil {
		return policy.EventPattern{}, fmt.Errorf("error parsing performs clause: %w", err)
	}

	containers := make([]string, 0)
	if index < len(fields) {
		on := fields[index+1:]
		containers = e.parseOn(strings.Join(on, " "))
	}

	return policy.EventPattern{
		Subject:    policy.Subject(subject),
		Operations: ops,
		Containers: containers,
	}, nil
}

func (e eventParser) parsePerforms(performs string) ([]policy.EventOperation, error) {
	ops := make([]policy.EventOperation, 0)
	split := strings.Split(performs, Or)
	hasArgs := false
	for _, s := range split {
		op := e.parseEventOperation(s)
		if len(op.Args) > 0 {
			hasArgs = true
		}

		ops = append(ops, op)
	}

	if hasArgs && len(ops) > 1 {
		return nil, fmt.Errorf("PERFORMS clause cannot have multiple operations if any of the operations have args")
	}

	return ops, nil
}

func (e eventParser) parseEventOperation(opStr string) policy.EventOperation {
	op := policy.EventOperation{}

	// if the string contains a parenthesis the operation has arguments
	if strings.Contains(opStr, "(") {
		split := strings.FieldsFunc(opStr, func(r rune) bool {
			return r == '(' || r == ')'
		})

		op.Operation = split[0]
		op.Args = strings.Fields(strings.ReplaceAll(split[1], ",", " "))
	} else {
		op.Operation = opStr
	}

	return op
}

func (e eventParser) parseOn(on string) []string {
	return strings.Fields(strings.ReplaceAll(on, ",", " "))
}

func (r responseParser) Parse(response string) (policy.ResponsePattern, error) {
	response = response[strings.Index(response, "(")+1 : strings.LastIndex(response, ")")]
	statements, _, err := Parse(response)
	if err != nil {
		return policy.ResponsePattern{}, err
	}

	return policy.ResponsePattern{
		Actions: statements,
	}, nil
}

package ngac

import "encoding/json"

type (
	Obligations interface {
		Add(obligation Obligation) error
		Remove(label string) error
		Get(label string) (Obligation, error)
		All() ([]Obligation, error)

		json.Marshaler
		json.Unmarshaler
	}

	Obligation struct {
		User     string
		Label    string
		Event    EventPattern
		Response ResponsePattern
	}

	EventPattern struct {
		Subject    string
		Operations []EventOperation
		Containers []string
	}

	EventOperation struct {
		Operation string
		Args      []string
	}

	ResponsePattern struct {
		Actions []Statement
	}
)

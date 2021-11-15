package ngac

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/ngac/graph"
)

type (
	Prohibitions interface {
		Add(prohibition Prohibition) error
		Get(subject string) ([]Prohibition, error)
		Delete(subject string, prohibitionName string) error

		json.Marshaler
		json.Unmarshaler
	}

	Prohibition struct {
		Name         string           `json:"name,omitempty"`
		Subject      string           `json:"subject,omitempty"`
		Containers   map[string]bool  `json:"containers,omitempty"`
		Operations   graph.Operations `json:"operations,omitempty"`
		Intersection bool             `json:"intersection,omitempty"`
	}
)

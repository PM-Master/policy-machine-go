package ngac

import (
	"encoding/json"
)

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
		User     string          `json:"user"`
		Label    string          `json:"label"`
		Event    EventPattern    `json:"event"`
		Response ResponsePattern `json:"response"`
	}

	EventPattern struct {
		Subject    string           `json:"subject"`
		Operations []EventOperation `json:"operations"`
		Containers []string         `json:"containers"`
	}

	EventOperation struct {
		Operation string   `json:"operation"`
		Args      []string `json:"args"`
	}

	ResponsePattern struct {
		Actions []Statement `json:"actions"`
	}

	jsonObligation struct {
		User     string       `json:"user"`
		Label    string       `json:"label"`
		Event    EventPattern `json:"event"`
		Response jsonResponse `json:"response"`
	}

	jsonResponse struct {
		Actions []map[string][]byte
	}
)

func (o *Obligation) MarshalJSON() ([]byte, error) {
	actions := make([]map[string][]byte, 0)
	for _, action := range o.Response.Actions {
		var actionName string
		switch action.(type) {
		case *CreatePolicyStatement:
			actionName = "CreatePolicyStatement"
		case *CreateNodeStatement:
			actionName = "CreateNodeStatement"
		case *AssignStatement:
			actionName = "AssignStatement"
		case *DeassignStatement:
			actionName = "DeassignStatement"
		case *DeleteNodeStatement:
			actionName = "DeleteNodeStatement"
		case *GrantStatement:
			actionName = "GrantStatement"
		case *DenyStatement:
			actionName = "DenyStatement"
		case *ObligationStatement:
			actionName = "ObligationStatement"
		}

		bytes, err := json.Marshal(action)
		if err != nil {
			return nil, err
		}

		actions = append(actions, map[string][]byte{actionName: bytes})
	}

	return json.Marshal(jsonObligation{
		User:     o.User,
		Label:    o.Label,
		Event:    o.Event,
		Response: jsonResponse{actions},
	})
}

func (o *Obligation) UnmarshalJSON(bytes []byte) error {
	j := jsonObligation{}
	err := json.Unmarshal(bytes, &j)
	if err != nil {
		return err
	}

	actions := make([]Statement, 0)

	for _, actionMap := range j.Response.Actions {
		for actionType, actionBytes := range actionMap {
			switch actionType {
			case "CreatePolicyStatement":
				action := CreatePolicyStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "CreateNodeStatement":
				action := CreateNodeStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "AssignStatement":
				action := AssignStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "DeassignStatement":
				action := DeassignStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "DeleteNodeStatement":
				action := DeleteNodeStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "GrantStatement":
				action := GrantStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "DenyStatement":
				action := DenyStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			case "ObligationStatement":
				action := ObligationStatement{}
				err = json.Unmarshal(actionBytes, &action)
				if err != nil {
					return err
				}
				actions = append(actions, &action)
			}
		}
	}

	o.Label = j.Label
	o.User = j.User
	o.Event = j.Event
	o.Response = ResponsePattern{Actions: actions}

	return nil
}

package memory

import (
	"encoding/json"
	"github.com/PM-Master/policy-machine-go/ngac"
)

type prohibitions struct {
	prohibitions map[string][]ngac.Prohibition
}

func NewProhibitions() ngac.Prohibitions {
	return &prohibitions{
		prohibitions: make(map[string][]ngac.Prohibition),
	}
}

func (p *prohibitions) Add(prohibition ngac.Prohibition) error {
	var (
		subjectPros []ngac.Prohibition
		ok          bool
	)

	if subjectPros, ok = p.prohibitions[prohibition.Subject]; !ok {
		subjectPros = make([]ngac.Prohibition, 0)
	}

	subjectPros = append(subjectPros, prohibition)
	p.prohibitions[prohibition.Subject] = subjectPros

	return nil
}

func (p *prohibitions) Get(subject string) ([]ngac.Prohibition, error) {
	return p.prohibitions[subject], nil
}

func (p *prohibitions) Delete(subject string, prohibitionName string) error {
	subjectPros := p.prohibitions[subject]
	newSubjectPros := make([]ngac.Prohibition, 0)
	for _, p := range subjectPros {
		if p.Name == prohibitionName {
			continue
		}

		newSubjectPros = append(newSubjectPros, p)
	}

	p.prohibitions[subject] = newSubjectPros
	return nil
}

type jsonProhibitions struct {
	Prohibitions []ngac.Prohibition `json:"prohibitions"`
}

func (p *prohibitions) MarshalJSON() ([]byte, error) {
	jp := jsonProhibitions{
		Prohibitions: make([]ngac.Prohibition, 0),
	}

	for _, subjectProhibitions := range p.prohibitions {
		for _, prohibition := range subjectProhibitions {
			jp.Prohibitions = append(jp.Prohibitions, prohibition)
		}
	}

	return json.Marshal(jp)
}

func (p *prohibitions) UnmarshalJSON(bytes []byte) error {
	jp := jsonProhibitions{
		Prohibitions: make([]ngac.Prohibition, 0),
	}

	if err := json.Unmarshal(bytes, &jp); err != nil {
		return err
	}

	p.prohibitions = make(map[string][]ngac.Prohibition)

	for _, prohibition := range jp.Prohibitions {
		subjectPros, ok := p.prohibitions[prohibition.Subject]
		if !ok {
			subjectPros = make([]ngac.Prohibition, 0)
		}

		subjectPros = append(subjectPros, prohibition)

		p.prohibitions[prohibition.Subject] = subjectPros
	}

	return nil
}

package policy

type (
	Prohibition struct {
		Name         string          `json:"name,omitempty"`
		Subject      string          `json:"subject,omitempty"`
		Containers   map[string]bool `json:"containers,omitempty"`
		Operations   Operations      `json:"operations,omitempty"`
		Intersection bool            `json:"intersection,omitempty"`
	}
)

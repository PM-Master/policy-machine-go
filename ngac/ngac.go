// Package ngac contains interfaces and structs necessary to be compliant with the NGAC standard
package ngac

type (
	FunctionalEntity interface {
		Graph() Graph
		Prohibitions() Prohibitions
		Obligations() Obligations
	}
)

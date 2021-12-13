package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/policy"
	"io/ioutil"
	"os"
)

type Properties map[string]string

func ReadAndApply(policyStore policy.Store, path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("recevied directory path, expected .ngac file")
	}

	var pal []byte
	if pal, err = ioutil.ReadFile(path); err != nil {
		return fmt.Errorf("error reading file %q: %w", fileInfo.Name(), err)
	}

	return apply(policyStore, string(pal))
}

func apply(policyStore policy.Store, pal string) error {
	stmts, _, err := Parse(pal)
	if err != nil {
		return fmt.Errorf("error parsing policy author language: %w", err)
	}

	for _, stmt := range stmts {
		err = stmt.Apply(policyStore)
		if err != nil {
			return fmt.Errorf("error applying statement: %w", err)
		}
	}

	return nil
}

func Author(policyStore policy.Store, stmts ...policy.Statement) error {
	for i, stmt := range stmts {
		if err := stmt.Apply(policyStore); err != nil {
			return fmt.Errorf("error applying statement at index %v: %v", i, err)
		}
	}

	return nil
}

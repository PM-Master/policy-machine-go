package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"io/ioutil"
	"os"
)

type Author struct {
	fe  ngac.FunctionalEntity
	pal string
}

func New(fe ngac.FunctionalEntity) Author {
	return Author{fe: fe}
}

func (a *Author) ReadAndApply(path string) error {
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

	a.pal = string(pal)

	return a.apply()
}

func (a Author) apply() error {
	stmts, _, err := Parse(a.pal)
	if err != nil {
		return fmt.Errorf("error parsing policy author language: %w", err)
	}

	for _, stmt := range stmts {
		err = stmt.Apply(a.fe)
		if err != nil {
			return fmt.Errorf("error applying statement: %w", err)
		}
	}

	return nil
}

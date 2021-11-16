package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type (
	Function struct {
		Name string
		Args []string
		Exec func(fe ngac.FunctionalEntity) error
	}

	Author struct {
		fe        ngac.FunctionalEntity
		pals      []string
		functions map[string]ParsedFunction
	}
)

func New(fe ngac.FunctionalEntity) Author {
	return Author{fe: fe}
}

func (a *Author) ReadPAL(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	files := make([]string, 0)

	if fileInfo.IsDir() {
		if err := filepath.WalkDir(fileInfo.Name(), func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			files = append(files, path)
			return nil
		}); err != nil {
			return fmt.Errorf("error loading PAL files: %w", err)
		}
	} else {
		files = append(files, path)
	}

	functions := make(map[string]ParsedFunction)
	pals := make([]string, 0)
	for _, fileName := range files {
		var pal []byte
		if pal, err = ioutil.ReadFile(fileName); err != nil {
			return fmt.Errorf("error reading file %q: %w", fileName, err)
		}

		pals = append(pals, string(pal))

		// parse functions
		_, parsedFunctions, err := Parse(string(pal))
		if err != nil {
			return err
		}

		for _, function := range parsedFunctions {
			if _, ok := functions[function.Name]; ok {
				return fmt.Errorf("function with name %q already exists", function.Name)
			}

			functions[function.Name] = function
		}
	}

	a.pals = pals
	a.functions = functions

	return nil
}

func (a Author) Apply() error {
	for _, pal := range a.pals {
		stmts, _, err := Parse(pal)
		if err != nil {
			return fmt.Errorf("error parsing policy author language: %w", err)
		}

		for _, stmt := range stmts {
			err = stmt.Apply(a.fe)
			if err != nil {
				return fmt.Errorf("error applying statement: %w", err)
			}
		}
	}

	return nil
}

func (a Author) Exec(fe ngac.FunctionalEntity, funcName string, args map[string]string) error {
	function, ok := a.functions[funcName]
	if !ok {
		return fmt.Errorf("function %q does not exist", funcName)
	}

	if len(function.Args) != len(args) {
		return fmt.Errorf("expected %d args for function %q but recevied %d",
			len(function.Args), function.Name, len(args))
	}

	for arg := range args {
		if !function.Args[arg] {
			return fmt.Errorf("unknown arg %q", arg)
		}
	}

	stmtsStr := function.Stmts
	for argName, argValue := range args {
		argName = fmt.Sprintf("$%s", argName)
		stmtsStr = strings.ReplaceAll(stmtsStr, argName, argValue)
	}

	stmts, _, err := Parse(stmtsStr)
	if err != nil {
		return err
	}

	for _, stmt := range stmts {
		if err := stmt.Apply(fe); err != nil {
			return err
		}
	}

	return nil
}

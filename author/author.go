package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Author struct {
	fe   ngac.FunctionalEntity
	pals []string
}

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

	pals := make([]string, 0)
	for _, fileName := range files {
		var pal []byte
		if pal, err = ioutil.ReadFile(fileName); err != nil {
			return fmt.Errorf("error reading file %q: %w", fileName, err)
		}

		pals = append(pals, string(pal))
	}

	a.pals = pals

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

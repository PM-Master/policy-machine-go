package author

import (
	"fmt"
	"github.com/PM-Master/policy-machine-go/ngac"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	Author interface {
		Apply(fe ngac.FunctionalEntity) error
	}

	author struct {
		pals []string
	}
)

func New(path string) (Author, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("error loading PAL files: %w", err)
		}
	} else {
		files = append(files, fileInfo.Name())
	}

	pals := make([]string, 0)

	for _, fileName := range files {
		var pal []byte
		if pal, err = ioutil.ReadFile(fileName); err != nil {
			return nil, fmt.Errorf("error reading file %q: %w", fileName, err)
		}

		pals = append(pals, string(pal))
	}

	return author{pals}, nil
}

func (a author) Apply(fe ngac.FunctionalEntity) error {
	for _, pal := range a.pals {
		stmts, err := Parse(pal)
		if err != nil {
			return fmt.Errorf("error parsing policy author language: %w", err)
		}

		for _, stmt := range stmts {
			err = stmt.Apply(fe)
			if err != nil {
				return fmt.Errorf("error applying statement: %w", err)
			}
		}
	}

	return nil
}

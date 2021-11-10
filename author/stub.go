package author

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GenerateFunctionsStub generates a stub for all functions found in the input path
// Functions will be outputted in Go.
// Functions that start with lower case letters will have lower case letters in the generated stub.
func GenerateFunctionsStub(name string, inputPath string, outputPath string) error {
	fileInfo, err := os.Stat(inputPath)
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
		files = append(files, fileInfo.Name())
	}

	pals := make([]string, 0)
	for _, fileName := range files {
		var pal []byte
		if pal, err = ioutil.ReadFile(fileName); err != nil {
			return fmt.Errorf("error reading file %q: %w", fileName, err)
		}

		pals = append(pals, string(pal))
	}

	outputFile, err := os.OpenFile(outputPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	s := initStub(name)
	for _, pal := range pals {
		_, functions, err := Parse(pal)
		if err != nil {
			return fmt.Errorf("error parsing policy author language: %w", err)
		}

		for _, function := range functions {
			s = fmt.Sprintf("%s\n%s", s, generateFunctionStub(name, function))
		}
	}

	_, err = outputFile.WriteString(s)
	return err
}

func initStub(name string) string {
	s := `package %s

import (
	"github.com/PM-Master/policy-machine-go/author"
	"github.com/PM-Master/policy-machine-go/ngac"
)

type %s struct {
  author author.Author
  functionalEntity ngac.FunctionalEntity
}

func New%s (author author.Author, functionalEntity ngac.FunctionalEntity) %s {
	return %s{author: author, functionalEntity: functionalEntity}
}`

	return fmt.Sprintf(s, name, name, name, name, name)
}

func generateFunctionStub(name string, function ParsedFunction) string {
	paramStr := ""
	for arg := range function.Args {
		paramStr += arg + " string,"
	}
	paramStr = strings.TrimSuffix(paramStr, ",")

	argStr := "map[string]string{%s}"
	argSubStr := ""
	for argName := range function.Args {
		s := fmt.Sprintf("%q:%s", argName, argName)
		argSubStr += fmt.Sprintf("%s,", s)
	}
	argStr = fmt.Sprintf(argStr, argSubStr)

	return fmt.Sprintf(`
func (stub %s) %s(%s) error {
  return stub.author.Exec(stub.functionalEntity, %q, %s)
}
`, name, function.Name, paramStr, function.Name, argStr)
}

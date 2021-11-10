package testoutput

import (
	"github.com/PM-Master/policy-machine-go/author"
	"github.com/PM-Master/policy-machine-go/ngac"
)

type TestStub struct {
	author           author.Author
	functionalEntity ngac.FunctionalEntity
}

func NewTestStub(author author.Author, functionalEntity ngac.FunctionalEntity) TestStub {
	return TestStub{author: author, functionalEntity: functionalEntity}
}

func (stub TestStub) my_function3(arg1 string) error {
	return stub.author.Exec(stub.functionalEntity, "my_function3", map[string]string{"arg1": arg1})
}

func (stub TestStub) my_function1(arg1 string, arg2 string) error {
	return stub.author.Exec(stub.functionalEntity, "my_function1", map[string]string{"arg1": arg1, "arg2": arg2})
}

func (stub TestStub) my_function2() error {
	return stub.author.Exec(stub.functionalEntity, "my_function2", map[string]string{})
}

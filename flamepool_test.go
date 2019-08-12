package flamepool

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	pool := New(20, []string{"tom"})

	if pool == nil {
		t.Error("Expected initialized pool")
	}
}

func TestParamsAreValid(t *testing.T) {
	valid := paramsAreValid("func (string, string) ( string, error)", "tom")
	if valid == false {
		t.Error("Expected valid params")
	}

	valid = paramsAreValid("func () (string, error)", "tom", "tom")
	if valid == true {
		t.Error("Expected invalid params")
	}
}

func TestReturnsAreValid(t *testing.T) {
	valid := returnsAreValid("func (string, string) ( string, error)")
	if valid == false {
		t.Error("Expected valid params")
	}

	valid = returnsAreValid("func (string, string)  (string)")
	if valid == true {
		t.Error("Expected invalid params")
	}
}

func TestRunFnSuccess(t *testing.T) {
	pool := New(1, []string{"tom"})
	fn := func(element string, othervar string) (interface{}, error) {
		return element + "--", nil
	}
	pool.Run(fn, "test")
}

func TestRunFnWithErr(t *testing.T) {
	pool := New(1, []string{"tom", "b"})
	fn := func(element string) (interface{}, error) {
		if element == "b" {
			return nil, errors.New("error")
		}
		return element + "--", nil
	}
	pool.Run(fn)
}

func TestRunTaskFailBecauseNoTaskReceived(t *testing.T) {
	pool := New(1, []string{"rulo", "tomcat"})
	NotATask := struct {
		Foo string
	}{
		"bar",
	}
	_, err := pool.Run(NotATask)
	if err == nil {
		t.Errorf("poo.Run() received a struct but isn't a Task")
	}
}

func TestErrorChannelReception(t *testing.T) {
	elements := []string{"rulo", "tomcat"}
	pool := New(1, elements)

	// This task has a Do() that always inject an error in the error channel
	task := FooTaskAlwaysDoError{}
	result, err := pool.Run(task)

	if err != nil {
		t.Errorf("Run failed")
	}

	if len(result.Errors) != len(elements) {
		errorDetail := fmt.Sprint("Mismatched errors quantity. Expected: ", elements, " Received: ", result.Errors)
		t.Errorf(errorDetail)
	}

	if len(result.Successful) != 0 {
		t.Errorf("It's supposed that there's no messages on successful results channel")
	}
}

func TestResultChannelReception(t *testing.T) {
	elements := []string{"rulo", "tomcat"}
	pool := New(1, elements)

	// This task has a Do() that always inject an successful in the result channel
	task := FooTaskAlwaysDoSuccessfulResult{}
	result, err := pool.Run(task)

	if err != nil {
		t.Errorf("Run failed")
	}

	if len(result.Errors) != 0 {
		t.Errorf("It's supposed that there's no messages on error channel")
	}

	if len(result.Successful) != len(elements) {
		errorDetail := fmt.Sprint("Mismatched results quantity. Expected: ", elements, " Received: ", result.Successful)
		t.Errorf(errorDetail)
	}
}

type FooTaskAlwaysDoError struct{}

func (ft FooTaskAlwaysDoError) Do(element interface{}) (interface{}, error) {
	value := reflect.ValueOf(element).Interface().(string)
	return nil, errors.New(value)
}

type FooTaskAlwaysDoSuccessfulResult struct{}

func (ft FooTaskAlwaysDoSuccessfulResult) Do(element interface{}) (interface{}, error) {
	value := reflect.ValueOf(element).Interface().(string)
	return value, nil
}

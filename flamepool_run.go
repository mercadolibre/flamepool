package flamepool

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// FlameResults response
type FlameResults struct {
	Successful []interface{}
	Errors     []error
}

// Task Flamepool
type Task interface {
	Do(element interface{}) (interface{}, error)
}

func (pool *Pool) run(obj interface{}, args ...interface{}) (FlameResults, error) {
	v := reflect.ValueOf(obj)

	switch v.Kind() {
	case reflect.Func:
		fnStr := fmt.Sprint(reflect.TypeOf(obj))
		if !isValid(fnStr, args...) {
			return FlameResults{}, errors.New("unexpected number of params, one param for each element required")
		}

		rargs := make([]reflect.Value, len(args)+1)
		for i, a := range args {
			rargs[i+1] = reflect.ValueOf(a)
		}

		v := reflect.ValueOf(obj)
		for i := 0; i < pool.poolSize; i++ {
			go pool.turnOnFnWorkers(v, rargs)
		}

	case reflect.Struct:
		task, ok := obj.(Task)
		if !ok {
			return FlameResults{}, errors.New("the struct must be a task")
		}

		for i := 0; i < pool.poolSize; i++ {
			go pool.turnOnTaskWorkers(task)
		}

	default:
		return FlameResults{}, errors.New("invalid type")
	}

	for _, element := range pool.elements {
		pool.innerChan <- element
	}

	var results []interface{}
	var errors []error

	for range pool.elements {
		select {
		case result := <-pool.resultChan:
			results = append(results, result)
		case err := <-pool.errorChan:
			errors = append(errors, err)
		}
	}

	flameResult := FlameResults{
		Successful: results,
		Errors:     errors,
	}

	return flameResult, nil
}

// isValid checks whether function parameters and returns are valid
func isValid(kind string, args ...interface{}) bool {
	splitted := strings.Split(kind, "(")
	if len(splitted) < 3 {
		return false
	}

	params := strings.Split(splitted[1], ",")
	returns := strings.Split(splitted[2], ",")

	if len(params) == len(args)+1 && len(returns) == 2 {
		return true
	}

	return false
}

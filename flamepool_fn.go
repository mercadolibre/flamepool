package flamepool

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (pool *Pool) runFn(fn interface{}, args ...interface{}) (FlameResults, error) {
	flameresult := FlameResults{}

	fnStr := fmt.Sprint(reflect.TypeOf(fn))
	if !paramsAreValid(fnStr, args...) {
		return flameresult, errors.New("Params invalid, one param for each element required")
	} else if !returnsAreValid(fnStr) {
		return flameresult, errors.New("Number of params invalid")
	}
	rargs := make([]reflect.Value, len(args)+1)
	for i, a := range args {
		rargs[i+1] = reflect.ValueOf(a)
	}

	v := reflect.ValueOf(fn)
	for i := 0; i < pool.poolSize; i++ {
		go pool.turnOnFnWorkers(v, rargs)
	}

	for _, element := range pool.elements {
		pool.innerChan <- element
	}

	results := []interface{}{}
	errors := []error{}
	for range pool.elements {
		select {
		case result := <-pool.resultChan:
			results = append(results, result)
		case err := <-pool.errorChan:
			errors = append(errors, err)
		}
	}
	flameresult.Successful = results
	flameresult.Errors = errors

	return flameresult, nil
}

func (pool *Pool) turnOnFnWorkers(fn reflect.Value, params []reflect.Value) {
	for element := range pool.innerChan {
		paramsNew := make([]reflect.Value, len(params))
		copy(paramsNew, params)
		paramsNew[0] = reflect.ValueOf(element)

		results := fn.Call(paramsNew)
		res, err := results[0].Interface(), results[1].Interface()
		if err != nil {
			pool.errorChan <- err.(error)
		} else {
			pool.resultChan <- res
		}
	}
}

func paramsAreValid(kind string, args ...interface{}) bool {
	splitted := strings.Split(kind, "(")
	params := strings.Split(splitted[1], ",")

	if len(params) == len(args)+1 {
		return true
	}
	return false
}

func returnsAreValid(kind string) bool {
	splitted := strings.Split(kind, "(")
	if len(splitted) < 3 {
		return false
	}
	params := strings.Split(splitted[2], ",")

	if len(params) == 2 {
		return true
	}
	return false
}

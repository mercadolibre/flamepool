package flamepool

import (
	"reflect"
)

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

func (pool *Pool) turnOnTaskWorkers(task Task) {
	for element := range pool.innerChan {
		res, err := task.Do(element)
		if err != nil {
			pool.errorChan <- err
		} else {
			pool.resultChan <- res
		}
	}
}

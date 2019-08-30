package flamepool

import (
	"errors"
)

// Task Flamepool
type Task interface {
	Do(element interface{}) (interface{}, error)
}

func (pool *Pool) runTask(obj interface{}) (FlameResults, error) {
	flameresult := FlameResults{}

	task, ok := obj.(Task)
	if !ok {
		return flameresult, errors.New("the struct must be a task")
	}

	for i := 0; i < pool.poolSize; i++ {
		go pool.turnOnTaskWorkers(task)
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

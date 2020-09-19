package flamepool

import (
	"fmt"
	"reflect"
)

// Pool Task
type Pool struct {
	poolSize   int
	resultChan chan interface{}
	errorChan  chan error
	innerChan  chan interface{}
	elements   []interface{}

	fn     reflect.Value
	params []reflect.Value
}

// New pool task
func New(poolSize int, items interface{}) *Pool {
	return newPool(poolSize, items)
}

func newPool(poolSize int, items interface{}) *Pool {
	var elements []interface{}

	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			elements = append(elements, elem)
		}
	}

	pool := &Pool{
		poolSize:   poolSize,
		elements:   elements,
		resultChan: make(chan interface{}, len(elements)),
		errorChan:  make(chan error, len(elements)),
		innerChan:  make(chan interface{}, len(elements)),
	}

	return pool
}

// Run task
func (pool *Pool) Run(obj interface{}, args ...interface{}) (FlameResults, error) {
	flameResult, err := pool.run(obj, args...)
	if err != nil {
		return FlameResults{}, fmt.Errorf("error: %v", err)
	}

	return flameResult, nil
}

// ChangeSettings for pool
func (pool *Pool) ChangeSettings(poolSize int, items interface{}) {
	newPool := newPool(poolSize, items)

	*pool = *newPool
}

// Close Closes all channels and resources associated with pool
func (pool *Pool) Close() {
	close(pool.resultChan)
	close(pool.errorChan)
	close(pool.innerChan)
}

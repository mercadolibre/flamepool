package flamepool

import (
	"errors"
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
func New(poolSize int, items ...interface{}) *Pool {
	return newPool(poolSize, items)
}

func newPool(poolSize int, items []interface{}) *Pool {
	pool := &Pool{}
	pool.poolSize = poolSize

	elements := make([]interface{}, len(items))
	copy(elements, items)

	pool.elements = elements
	pool.resultChan = make(chan interface{}, len(pool.elements))
	pool.errorChan = make(chan error, len(pool.elements))
	pool.innerChan = make(chan interface{}, len(pool.elements))

	return pool
}

// Run task
func (pool *Pool) Run(obj interface{}, args ...interface{}) (FlameResults, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Func {
		return pool.runFn(obj, args...)
	} else if v.Kind() == reflect.Struct {
		return pool.runTask(obj)
	}
	return FlameResults{}, errors.New("invalid type")
}

// ChangeSettings for pool
func (pool *Pool) ChangeSettings(poolSize int, items ...interface{}) {
	newPool := newPool(poolSize, items)

	*pool = *newPool
}

// Close Closes all channels and resources associated with pool
func (pool *Pool) Close() {
	close(pool.resultChan)
	close(pool.errorChan)
	close(pool.innerChan)
}

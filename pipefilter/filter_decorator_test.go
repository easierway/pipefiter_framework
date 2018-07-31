package pipefilter

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/easierway/service_decorators"
)

type demoFilter struct {
	executionCnt int
}

func (f *demoFilter) Process(data interface{}) (interface{}, error) {
	var (
		str string
		ok  bool
	)
	if str, ok = data.(string); !ok {
		panic("invalid input")
	}
	fmt.Println("input:", str)
	f.executionCnt++
	if f.executionCnt < 3 {
		return nil, errors.New("error occurred")
	}
	return f.executionCnt, nil
}

func TestDecoratorFilter(t *testing.T) {
	var (
		retryDec *service_decorators.RetryDecorator
		err      error
	)
	retriableChecker := func(err error) bool {
		return true
	}

	if retryDec, err = service_decorators.CreateRetryDecorator(5, /*max retry times*/
		time.Second*1, time.Second*2, retriableChecker); err != nil {
		panic(err)
	}
	decorators := []service_decorators.Decorator{retryDec}
	orgFilter := &demoFilter{}
	decoratedFilter := DecorateFilter(orgFilter, decorators)
	startT := time.Now()
	ret, err := decoratedFilter.Process("Hello")
	timeSpent := time.Now().Sub(startT).Nanoseconds()
	fmt.Println("time escaped: ", timeSpent)
	if err != nil {
		t.Error("Unexpected error occurred.")
	}
	if ret.(int) != 3 {
		t.Errorf("Expected value is %d, but actual value is %d", 3, ret)
	}
}

func ExampleDecoratedFilter() {
	// 1. Create the decorators
	var (
		retryDec *service_decorators.RetryDecorator
		err      error
	)
	retriableChecker := func(err error) bool {
		return true
	}

	if retryDec, err = service_decorators.CreateRetryDecorator(5, /*max retry times*/
		time.Second*1, time.Second*2, retriableChecker); err != nil {
		panic(err)
	}
	// Put all decorators in the slice.
	// Be careful of the order of the decorators, your filter will be decorated by the same order
	// and the decorators will be invoked by the same order
	decorators := []service_decorators.Decorator{retryDec}

	// 2. Decorate the filter with the decorators by calling DecorateFilter
	// after that you will get a deocrated Filter instance
	orgFilter := &demoFilter{}
	decoratedFilter := DecorateFilter(orgFilter, decorators)

	// 3. The decorated instance can be used as a normal Fiter instance
	decoratedFilter.Process("Hello")
}

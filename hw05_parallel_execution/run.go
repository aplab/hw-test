package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		m = len(tasks) + 1
	}
	var errCounter uint64
	taskCh := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if atomic.LoadUint64(&errCounter) >= uint64(m) {
					break
				}
				if task() != nil {
					atomic.AddUint64(&errCounter, 1)
				}
			}
		}()
	}
	for _, task := range tasks {
		if atomic.LoadUint64(&errCounter) >= uint64(m) {
			break
		}
		taskCh <- task
	}
	close(taskCh)
	wg.Wait()
	if errCounter >= uint64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

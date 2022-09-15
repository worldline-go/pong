package registry

import "sync"

type Errors struct {
	Errs  []error
	mutex sync.Mutex
}

func (e *Errors) AddError(err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Errs = append(e.Errs, err)
}

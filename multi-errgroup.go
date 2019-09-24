package errgroup

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

type MultiErrorGroup struct {
	wg     sync.WaitGroup
	errors []error
	lock   sync.Mutex
}

func (g *MultiErrorGroup) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.lock.Lock()
			defer g.lock.Unlock()

			g.errors = append(g.errors, err)
		}
	}()
}

func (g *MultiErrorGroup) Wait() error {
	g.wg.Wait()

	return multierror.Append(nil, g.errors...).ErrorOrNil()
}



package counter

import (
	"sync"
	"sync/atomic"
	"testing"
)

// -----------------------------------------------------------------------------

func TestCounter(t *testing.T) {
	var nbRoutines int = 8e3

	wg := &sync.WaitGroup{}
	wg.Add(nbRoutines)

	var total int64
	stop := make(chan struct{}, nbRoutines)
	c := NewCounter(0, 10, func(i int64) {
		atomic.AddInt64(&total, 1)
		stop <- struct{}{}
	})

	for i := 0; i < nbRoutines; i++ {
		go func() {
			for {
				select {
				case <-stop:
					wg.Done()
					return
				default:
					c()
				}
			}
		}()
	}

	wg.Wait()

	if total != int64(nbRoutines) {
		t.Errorf("expected total = %d, not %d", nbRoutines, total)
	}
	if next := c(); next != 1 {
		t.Errorf("expected c() = 1, not %d", next)
	}
}

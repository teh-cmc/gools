package errlog

import (
	"fmt"
	"sync"
	"testing"
)

// -----------------------------------------------------------------------------

func TestErrlog(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		for i := 0; i < 1e5; i++ {
			ErrLogger.Toggle(i%2 == 0)
		}
		wg.Done()
	}()
	go func() {
		err := fmt.Errorf("error")
		for i := 0; i < 1e5; i++ {
			ErrLogger.LogError(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

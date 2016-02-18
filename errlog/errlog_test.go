package errlog

import (
	"fmt"
	"sync"
	"testing"
)

// -----------------------------------------------------------------------------

func TestErrlog(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(3)
	go func() {
		for i := 0; i < 1e6; i++ {
			ErrLogger.Toggle(i%2 == 0)
		}
		wg.Done()
	}()
	go func() {
		err := fmt.Errorf("error")
		for i := 0; i < 1e6; i++ {
			_ = ErrLogger.Log(err)
		}
		wg.Done()
	}()
	go func() {
		err := fmt.Errorf("errorST")
		for i := 0; i < 1e6; i++ {
			_ = ErrLogger.LogST(err)
		}
		wg.Done()
	}()

	wg.Wait()
}

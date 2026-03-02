package keymux_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/keymux"
)

func TestMutex(t *testing.T) {
	t.Parallel()

	mux := keymux.New[string]()

	var group sync.WaitGroup
	group.Add(10000)
	for range 10000 {
		go assert.Eventually(t, func() bool {
			if mux.TryLock("test") {
				group.Done()
				mux.Unlock("test")
				return true
			}

			return false
		}, time.Second, time.Millisecond)
	}

	group.Wait()
}

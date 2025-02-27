package keylock

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestMemoryLock(t *testing.T) {
	var wg sync.WaitGroup

	n := 10
	cnt := 0
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			Memory().Lock("1")
			defer Memory().Unlock("1")
			cnt++
		}()
	}
	wg.Wait()

	require.Equal(t, n, cnt)

	ml := Memory().(*memoryLockerImpl)
	require.Empty(t, ml.locksMap)
}

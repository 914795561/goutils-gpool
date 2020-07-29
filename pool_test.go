package gpool

import (
	"sync"
	"sync/atomic"
	"testing"
)

var wg sync.WaitGroup
var sum int64
var runTimes = 1000000

func task(i int) func() {
	return func() {
		defer wg.Done()
		for a := 0; a < 100; a++ {
			atomic.AddInt64(&sum, int64(a))
		}
	}

}

func BenchmarkPool(b *testing.B) {
	p := NewPool(20)

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Add(task(i))
	}

	wg.Wait()
}

func BenchmarkGoroutine(b *testing.B) {

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go task(i)()
	}

	wg.Wait()
}
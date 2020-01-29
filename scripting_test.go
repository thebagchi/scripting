package scripting_test

import (
	"sync"
)

func parallel(count, threads int, operation func(id, thread int)) {
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		s, e := count/threads*i, count/threads*(i+1)
		if i == threads-1 {
			e = count
		}
		go func(i, s, e int) {
			defer wg.Done()
			for j := s; j < e; j++ {
				operation(j, i)
			}
		}(i, s, e)
	}
	wg.Wait()
}

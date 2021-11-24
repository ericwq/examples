package ptest

import "sync"

var (
	data   = make(map[string]string)
	locker sync.RWMutex
)

func WriteToMap(k, v string) {
	locker.Lock()
	defer locker.Unlock()
	data[k] = v
}

func ReadFromMap(k string) string {
	locker.RLock()
	defer locker.RUnlock()
	return data[k]
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func Sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

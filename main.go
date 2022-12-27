package main

import (
	"log"
	"sort"
	"sync"
	"time"
)

func dump(c chan int, a []int) {
	for _, i := range a {
		c <- i
	}
	close(c)
}

type MemStorage struct {
	sync.RWMutex
	items map[int]int
}

func (ms *MemStorage) Set(key int) {
	ms.Lock()
	defer ms.Unlock()

	ms.items[key] += 1
}

func merge(c1, c2 chan int) {
	m := MemStorage{items: make(map[int]int)}
	time.Sleep(1 * time.Second)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			select {
			case var1, ok1 := <-c1:
				time.Sleep(1 * time.Second)
				if !ok1 {
					wg.Done()
					return
				}
				m.Set(var1)
			default:
				wg.Done()
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case var2, ok2 := <-c2:
				time.Sleep(1 * time.Second)
				if !ok2 {
					wg.Done()
					return
				}
				m.Set(var2)
			default:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	a := make([]int, 0)
	for k, v := range m.items {
		if v == 1 {
			a = append(a, k)
		}
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	log.Println(a)
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	a1 := []int{1, 5, 7, 8, 24, 25, 50}
	a2 := []int{1, 3, 6, 8, 25, 3123, 32423, 312312, 3212311, 3212312, 3212313}
	go dump(c1, a1)
	go dump(c2, a2)
	merge(c1, c2)
}

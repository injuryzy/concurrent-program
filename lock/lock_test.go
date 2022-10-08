package lock

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// 数据竞争  main 和  协程
func TestDataCompetition(t *testing.T) {
	var data int
	go func() {
		data++
	}()
	if data == 0 {
		fmt.Println(data)
	}
}

//dail lock
// v1 锁住 v2 也锁住了
func TestDailLock(t *testing.T) {
	type value struct {
		mu    sync.Mutex
		value int
	}
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Printf("sum=%d", v1.value+v2.value)
	}
	a := value{value: 1}
	b := value{value: 1}
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

func TestLiveLock(t *testing.T) {

	//cond := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {

		}
	}()
}
func TestName(t *testing.T) {

	for t2 := range time.Tick(1 * time.Second) {
		fmt.Println(t2)
	}
}

//测量一个grouting的大小
func TestGc(t *testing.T) {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}
	var c <-chan interface{}

	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c
	}
	const numGroutines = 1e4
	wg.Add(numGroutines)
	//获取gc初始状态
	before := memConsumed()
	for i := numGroutines; i > 0; i-- {

		go noop()
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%3fkb", float64(after-before)/numGroutines/1000)
}

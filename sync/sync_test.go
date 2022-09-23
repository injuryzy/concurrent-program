package sync

import (
	"fmt"
	"math"
	"os"
	"sync"
	"testing"
	"text/tabwriter"
	"time"
)

//waitgroup
func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st groutinnue sleeping ...")
		time.Sleep(1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2st grouting sleeping ...")
		time.Sleep(2)
	}()
	wg.Wait()
	fmt.Println("all groutines complete")

}

//mutex

func TestMutex(t *testing.T) {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("incrementing:%d\n", count)
	}
	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("decrement:%d\n", count)
	}

	var arithmetic sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}
	arithmetic.Wait()
	fmt.Println(count)
}

func TestRWMetux(t *testing.T) {
	//写锁
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			time.Sleep(1)
			l.Unlock()
		}
	}
	//读锁
	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, metux, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(1 + count)
		begin := time.Now()
		go producer(&wg, metux)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(begin)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	var m sync.RWMutex
	fmt.Fprintf(tw, "Reader\tRwmutex\tmutex\n")
	for i := 0; i < 10; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(tw, "%d\t%v\t%v\n", count, test(count, &m, m.RLocker()), test(count, &m, &m))
	}
}

func TestCond(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	remove := func(delay time.Duration) {
		time.Sleep(delay)
		cond.L.Lock()
		queue = queue[1:]
		fmt.Println("remove from queue")
		cond.L.Unlock()
		cond.Broadcast()
	}

	for i := 0; i < 10; i++ {
		cond.L.Lock()
		//这里不用if是因为要重复判断
		for len(queue) == 2 {
			cond.Wait()
		}

		queue = append(queue, i)
		fmt.Println("add to queue")
		go remove(1 * time.Second)
		cond.L.Unlock()
	}
}

func TestCondBroatCast(t *testing.T) {
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var tempwg sync.WaitGroup
		tempwg.Add(1)
		go func() {
			tempwg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		tempwg.Wait()
	}

	var wg sync.WaitGroup
	wg.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("maximizing window")
		wg.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("display window")
		wg.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("mouse window")
		wg.Done()
	})
	button.Clicked.Broadcast()
	wg.Wait()
}

func TestOnce(t *testing.T) {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			once.Do(increment)
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

func TestPool(t *testing.T) {
	p := &sync.Pool{New: func() any {
		fmt.Println("create new interfance")
		return struct{}{}
	}}
	p.Get()
	get := p.Get()
	p.Put(get)
	p.Get()
}

func TestPool2(t *testing.T) {
	var numCalcsCreated int
	calcPoll :=&sync.Pool{New: func() any{
		numCalcsCreated+=1

		return numCalcsCreated
	}}
	calcPoll.Put(calcPoll.New())

	const numworks =1024*1024
	var wg sync.WaitGroup
	wg.Add(numworks)
	for i := 0; i < numworks; i++ {
		go func() {
			defer wg.Done()
			i2 := calcPoll.Get().(int)
			fmt.Println(i2)
			defer calcPoll.Put(i2)
		}()
	}
	wg.Wait()
	fmt.Println("num",numCalcsCreated)
	fmt.Println("num1")

}

func TestCount(t *testing.T) {
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Println(count)
	}
}




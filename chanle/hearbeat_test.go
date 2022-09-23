package chanle

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestHearbeat(t *testing.T) {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			defer close(heartbeat)
			defer close(results)
			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			SendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				select {
				case <-done:
					return
				case <-pulse:
					SendPulse()
				case results <- r:
					return
				}
			}
			for {
				select {
				case <-done:
					return
				case <-pulse:
					SendPulse()
				case r := <-workGen:
					sendResult(r)

				}
			}
		}()
		return heartbeat, results
	}
	do := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(do)
	})
	const timeout = 2 * time.Second

	work, result := doWork(do, timeout/2)
	for {
		select {
		case _, ok := <-work:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-result:
			if ok == false {
				return
			}
			fmt.Println("results ", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}

//请求并发复制处理
func TestRequest(t *testing.T) {
	doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
		start := time.Now()
		defer wg.Done()
		duration := time.Duration(1+rand.Intn(5)) * time.Second
		select {
		case <-done:

		case <-time.After(duration):
		}
		select {
		case <-done:

		case result <- id:
		}
		took := time.Since(start)
		/*if took < duration {
			took = duration
		}*/
		fmt.Println(id, " took ", took)
	}
	done := make(chan interface{})
	result := make(chan int)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}
	firstRequest := <-result
	close(done)
	wg.Wait()

	fmt.Println(firstRequest)
}

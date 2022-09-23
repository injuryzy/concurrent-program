package chanle

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	stringStream := make(chan string)
	go func() {
		if false {

			stringStream <- "hello"
		}
	}()
	fmt.Println(<-stringStream)
}

func TestChanle(t *testing.T) {
	ints := make(chan int)

	go func() {
		//defer close(ints)
		time.Sleep(2 * time.Second)
		for i := 0; i < 5; i++ {
			ints <- i
		}
	}()

	for {
		select {
		case <-ints:
			fmt.Println(<-ints)
		case <-time.After(time.Second):
			fmt.Println("close")
			goto t
		}
	}
t:
	fmt.Println("return")
}

func TestDowork(t *testing.T) {
	doworkd := func(do <-chan any, string <-chan string) {
		go func() {
			for {
				time.Sleep(time.Millisecond)
				fmt.Println("1")
			}
			for {
				select {
				case <-do:

					return
				case <-string:
					fmt.Println("string")
				}
			}
		}()
	}
	do := make(chan any)
	doworkd(do, nil)
	time.Sleep(time.Second)
	close(do)
	fmt.Println("finish")
}

func TestOrChanle(t *testing.T) {

	var or func(chanle ...<-chan any) <-chan any
	or = func(chanle ...<-chan any) <-chan any {
		switch len(chanle) {
		case 0:
			return nil
		case 1:
			fmt.Println(1)
			return chanle[0]
		}
		orDone := make(chan any)
		go func() {
			defer close(orDone)
			switch len(chanle) {
			case 2:
				select {
				case <-chanle[0]:
					fmt.Println(2)
				case <-chanle[1]:
					fmt.Println(3)
				}
			default:
				select {
				case <-chanle[0]:
				case <-chanle[1]:
				case <-chanle[2]:
					fmt.Println(4)
				case <-or(append(chanle[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	sign := func(duration time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(duration)
		}()
		return c
	}

	fmt.Println(<-or(sign(3*time.Second), sign(2*time.Second), sign(time.Second), sign(time.Second)))
}

func TestError(t *testing.T) {
	type result struct {
		err     error
		reponse *http.Response
	}
	checkStatus := func(url ...string) <-chan result {
		results := make(chan result)
		go func() {
			defer close(results)
			for _, s := range url {
				var result result
				get, err := http.Get(s)
				if err != nil {
					result.err = err
				}
				result.reponse = get
				select {
				case results <- result:

				}
			}
		}()
		return results
	}
	status := checkStatus("a")
	r := <-status
	fmt.Println(r.err, r.reponse)
}

func TestDo(t *testing.T) {
	do := func(done <-chan any) {
		go func() {
			select {
			case <-done:
				fmt.Println(1)
			}
		}()
		go func() {
			//time.Sleep(time.Second)
			select {
			case <-done:
				fmt.Println(2)

			}
		}()
	}
	anies := make(chan any)
	do(anies)
	anies <- 1
	time.Sleep(2 * time.Second)
}
func TestBrangeChanle(t *testing.T) {
	brange := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case mstream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = mstream
				case <-done:
					return
				}
				valStream <- stream
			}
		}()
		return valStream
	}
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for i := range brange(nil, genVals()) {
		fmt.Println(i)
	}
	context.TODO()
}

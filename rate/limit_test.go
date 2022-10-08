package rate

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

type Apiconnection struct {
	rateLimter *rate.Limiter
}

func Openfile() *Apiconnection {
	return &Apiconnection{
		rateLimter: rate.NewLimiter(rate.Limit(1), 1),
	}
}

func (a *Apiconnection) ReadFile(ctx context.Context) error {
	if err := a.rateLimter.Wait(ctx); err != nil {

		return err
	}
	return nil
}

func (a *Apiconnection) ResolverAddr(ctx context.Context) error {
	if err := a.rateLimter.Wait(ctx); err != nil {
		return err
	}
	return nil

}
func TestRate(t *testing.T) {
	deadline, _ := context.WithDeadline(context.Background(), <-time.After(3*time.Second))
	openfile := Openfile()
	for true {
		err := openfile.ReadFile(deadline)
		fmt.Println(1, err)
		err = openfile.ResolverAddr(deadline)
		fmt.Println(2, err)
	}
}

func TestNewStart(t *testing.T) {
	var values = []int{1, 2, 3, 54, 66, 10}
	done := make(chan any)
	defer close(done)
	startFunc, stream := StartFunc(done, values)
	for {
		select {
		case val := <-startFunc:
			fmt.Println(val)
		case val2 := <-stream:
			fmt.Println(val2)

		}
	}
}
func StartFunc(done <-chan any, vals []int) (heartbet, stream chan any) {
	heartbet = make(chan any)
	stream = make(chan any)
	pulse := time.Tick(time.Second)
	single := time.Tick(3 * time.Second)

	go func() {
	Loop:
		for {

			for _, val := range vals {
				select {
				case <-done:
					return
				case <-pulse:
					heartbet <- struct{}{}
					stream <- val
				case <-single:
					StartFunc(done, vals)
					continue Loop
				}
			}
		}
	}()
	return heartbet, stream
}


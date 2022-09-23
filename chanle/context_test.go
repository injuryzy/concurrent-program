package chanle

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCancle(t *testing.T) {
	cancel, cfun := context.WithCancel(context.Background())
	defer cfun()
	go func() {
		select {
		case <-cancel.Done():
			fmt.Println(3241)
			fmt.Println(cancel.Err())
		}
	}()
	time.Sleep(time.Second * 2)
	fmt.Println(34312)
}
func TestTimeOut(t *testing.T) {
	cancel, cfun := context.WithTimeout(context.Background(), time.Second)
	defer cfun()
	go func() {
		select {
		case <-cancel.Done():
			fmt.Println(3241)
			fmt.Println(cancel.Err())
		}
	}()
	time.Sleep(time.Second * 2)
	fmt.Println(34312)
}

func TestWithDeadline(t *testing.T) {
	cancel, _ := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	value := context.WithValue(cancel, cancel, "value")
	//defer cfun()
	go func() {
		select {
		case <-cancel.Done():
			fmt.Println(3241)
			fmt.Println(cancel.Err())
		}
	}()
	time.Sleep(time.Second * 3)
	fmt.Println(34312)
	fmt.Println(value.Value(cancel))
}

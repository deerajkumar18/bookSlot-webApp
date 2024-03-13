package concurrency_test

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"testing"
)

type Res struct {
	resp *http.Response
	err  error
}

func PingBookEP(wg *sync.WaitGroup, cs chan Res, id int) {
	url := "http://localhost:8081/Book/1/" + strconv.Itoa(id)
	resp, err := http.Get(url)
	cs <- Res{resp, err}
	wg.Done()
}

func TestConcurrency_SlotBooking(t *testing.T) {

	wg := &sync.WaitGroup{}
	numberOfUsers := 10
	cs := make(chan Res)
	for i := 1; i <= numberOfUsers; i++ {
		wg.Add(1)
		go PingBookEP(wg, cs, i)
	}
	wg.Wait()
	close(cs)

	for i := 0; i < numberOfUsers; i++ {
		results := <-cs
		fmt.Println(i, results)
	}

}

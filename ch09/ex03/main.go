package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	cancel := make(chan struct{})
	memo := New(httpGetBody, cancel)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := memo.Get("http://gopl.io", cancel)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Println(res)
	}()
	close(cancel)
	wg.Wait()
}

type Func func(key string, cancel <-chan struct{}) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	response chan<- result
}

type Memo struct {
	requests chan request
}

func New(f Func, cancel <-chan struct{}) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f, cancel)
	return memo
}

func (memo *Memo) Get(key string, cancel <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	select {
	case <- cancel:
		return nil, errors.New(key + ": request is cancelled.")
	case res := <- response:
		return res.value, res.err
	}
}

func (memo *Memo) server(f Func, cancel <-chan struct{}) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go func() {
				e.call(f, req.key, cancel)
				select {
				case <-e.ready:
					close(e.ready)
				case <-cancel:
					fmt.Printf("%s: request cancelling...", req.key)
					for range memo.requests {
					}
					delete(cache, req.key)
					fmt.Printf("%s: request cancelled and delete cache.", req.key)
				}
			}()
		}
		go e.deliver(req.response, cancel)
	}
}

func (e *entry) call(f Func, key string, cancel <-chan struct{}) {
	fmt.Printf("%s: new request...\n", key)
	e.res.value, e.res.err = f(key, cancel)
	select {
	case <-cancel:
	default:
	}
	fmt.Printf("%s: get response...\n", key)
}

func (e *entry) deliver(response chan<- result, cancel <-chan struct{}) {
	select {
	case <-e.ready:
		response <- e.res
	case <-cancel:
		fmt.Println("deliver cancelled.")
	}
}

func httpGetBody(url string, cancel <-chan struct{}) (interface{}, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Cancel = cancel
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
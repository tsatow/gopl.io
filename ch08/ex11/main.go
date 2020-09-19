package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	response := mirroredQuery()
	fmt.Println(response)
}

func mirroredQuery() string {
	// 焦りからか非常にきたない
	responses := make(chan string)
	cancel := make(chan struct{})
	go func() {
		r, err := request("asia.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("europe.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("americas.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			return
		}
		responses <- r
	}()
	select {
	case r := <-responses:
		cancel <- struct{}{}
		return r
	}
}

func request(hostname string, cancel chan struct{}) (response string, err error) {
	fmt.Fprintf(os.Stdout, "request for %s\n", hostname)
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s", hostname), nil)
	req.Cancel = cancel
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "request is cancelled.")
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("getting %s: %s", hostname, resp.Status)
		return
	}
	// 無視
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("body is %s", b)
	response = fmt.Sprintf("%s", b)
	return
}

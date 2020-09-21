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
	os.Exit(0)
}

func mirroredQuery() string {
	// 焦りからか非常にきたない
	responses := make(chan string)
	cancel := make(chan struct{})
	go func() {
		r, err := request("gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("asia.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("europe.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		responses <- r
	}()
	go func() {
		r, err := request("americas.gopl.io", cancel)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		responses <- r
	}()
	select {
	case r := <-responses:
		go func() {
			// ここ、なんで詰まるんだっけ？
			cancel <- struct{}{}
			fmt.Println("cancelled.")
		}()
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
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	// 無視
	b, err := ioutil.ReadAll(resp.Body)
	response = fmt.Sprintf("%s", b)
	return
}

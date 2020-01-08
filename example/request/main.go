package main

import (
	"fmt"

	"github.com/fogfish/gurl"
)

type headers struct {
	UserAgent string `json:"X-User-Agent,omitempty"`
}

type httpbin struct {
	URL     string  `json:"url,omitempty"`
	Origin  string  `json:"origin,omitempty"`
	Headers headers `json:"headers,omitempty"`
}

func request() (val httpbin, err error) {
	err = gurl.IO().
		GET("https://httpbin.org/get").
		With("Accept", "application/json").
		With("X-User-Agent", "gurl").
		Code(200).
		Head("Content-Type", "application/json").
		Recv(&val).
		Fail

	return
}

func main() {
	val, err := request()

	if err != nil {
		fmt.Printf("fail %v\n", err)
	}
	fmt.Printf("==> %v\n", val)
}

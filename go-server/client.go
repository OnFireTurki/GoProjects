package main

import (
	"fmt"
	"io"
	"net/http"
)

func sendSimpleGet(url string) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ioR, _ := io.ReadAll(r.Body)
	fmt.Println(string(ioR))
}

func client() {
	sendSimpleGet("http://localhost:8080/hello") // Should return "Hello !!"
	sendSimpleGet("http://localhost:8080/form")  // Should return "405 method not allowed"
}

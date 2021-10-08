package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func sendMessegeError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	message := fmt.Sprintf(msg+"\n", args...)
	params := url.Values{}
	params.Add("erro", message)
	_, err := http.Get("<your-weservice>" + params.Encode())
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(1)
}

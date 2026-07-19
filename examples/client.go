// Copyright (c) BeduSec. All rights reserved.
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}
	url := "http://localhost:8080/"

	for i := 0; i < 5; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("Request %d: status=%d, body=%s\n", i+1, resp.StatusCode, body)
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Burst test: status=%d, body=%s\n", resp.StatusCode, body)

	time.Sleep(2 * time.Second)
	resp2, _ := client.Get(url)
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("After wait: status=%d, body=%s\n", resp2.StatusCode, body2)
}
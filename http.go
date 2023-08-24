package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpRequest(url string, headerAuthValue string) (data []byte) {
	clien := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Add("Authorization", headerAuthValue)
	// req.Header.Add("ContentType", headerContentType)

	resp, err := clien.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	// defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading HTTP response body:", err)
		return
	}

	// data = string(body)
	// fmt.Println(resp.StatusCode)
	return data

}

package main

import (
	"fmt"
	"time"
)

func displayLoading(loadingCh chan bool) {
	fmt.Print("Loading")
	for {
		select {
		case <-loadingCh:
			return
		default:
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func dispayJob(job string, phase string) {
	switch phase {
	case "start":
		fmt.Println("Start", job)
	case "end":
		fmt.Println("Finish", job)
	}
}

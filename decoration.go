package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func displayLoading(loadingCh chan bool) {
	fmt.Print("Loading")
	for {
		select {
		case <-loadingCh:
			fmt.Print()
			return
		default:
			fmt.Print(".")
			time.Sleep(100 * time.Millisecond) // Adjust sleep time as needed
		}
	}
}

func displayJob(job string, phase string) {
	switch phase {
	case "start":
		// clearScreen()
		fmt.Println("Start", job)
	case "end":
		fmt.Println("\nFinish", job)
		clearScreen()
	}
}

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin", "linux":
		cmd = exec.Command("clear")
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

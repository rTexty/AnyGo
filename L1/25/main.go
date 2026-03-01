package main

import (
	"fmt"
	"time"
)

func timeSleep(d time.Duration) {
    <-time.After(d)
}

func main() {
    start_time := time.Now()
    fmt.Println("Start", start_time)
    timeSleep(time.Duration(2 * time.Second))
    end_time := time.Now()
    fmt.Println("Finished", end_time)
    fmt.Println("Duration is: ", time.Since(start_time))
}
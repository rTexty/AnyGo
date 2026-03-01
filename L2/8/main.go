package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func getTime(address string) (time.Time,error) {
    t, err := ntp.Time(address)
    return t, err
}

func main() {
    time, err := getTime("0.beevik-ntp.pool.ntp.org")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error with ntpTime", err.Error())
        os.Exit(1)
    }
    fmt.Println(time)
}

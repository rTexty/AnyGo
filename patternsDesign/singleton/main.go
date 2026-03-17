package main

import (
	"fmt"
	"sync"
)

type Logger struct {
}

var once sync.Once
var logger *Logger

func getLogger() *Logger {
    if logger == nil{
        once.Do(func() {
            logger = &Logger{}
            fmt.Println("Logger initialization")
        })
    } else {
        fmt.Println("Logger already initialized")
        }
    return logger
}

func main() {
    for range 5{
        go getLogger()
    }
    fmt.Scanln()

}
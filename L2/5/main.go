package main

import "fmt"

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
//     println("error")
    fmt.Println()
    return
  }
  println("ok")
}
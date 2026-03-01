package main

import (
	"fmt"
	"strconv"
)

func main() {
  a := [5]int{76, 77, 78, 79, 80}
  var b []int = a[1:4]
  fmt.Println(b)
  fmt.Println([]byte("a4bc2d5e"))
  num, err:= strconv.Atoi("a")
  if err != nil {
      fmt.Println(err.Error())
  }
  fmt.Println(num)
}
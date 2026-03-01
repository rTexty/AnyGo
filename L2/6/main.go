package main

import (
	"fmt"
)

func main() {
  var s = make([]string,0, 10)
  s = append(s, []string{"1","2","3"}...)
  fmt.Println(len(s), cap(s))
  fmt.Println(s)
  modifySlice(s)
  fmt.Println(s)
}

func modifySlice(i []string) {
  i[0] = "3"
  i = append(i, "4")
  i[1] = "5"
  i = append(i, "6")
}
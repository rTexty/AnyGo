package main

import (
	"fmt"
	"strings"
)

func Validation(strscan string) bool{
    str := []rune(strings.ToLower(strscan))
    mapa := make(map[rune]int)
    for _,symb := range str {
        mapa[symb] += 1
        if mapa[symb] > 1{
            return false
        }
    }
    return true
}

func main() {

    fmt.Println("Введите строку: ")
    var strscan string
    fmt.Scan(&strscan)

    res := Validation(strscan)
    fmt.Println(res, strscan)
}
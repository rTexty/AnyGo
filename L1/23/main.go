package main

import "fmt"


func rem(sl []int, elem int) []int {
	copy(sl[elem:], sl[elem+1:])
    sl[len(sl)-1] = 0
	sl = sl[:len(sl)-1]
	return sl
}

func main() {
	sl1 := []int{1,2,3,4,5}
    elem := 2
	sl1 = rem(sl1, elem)
	fmt.Println(sl1)
    fmt.Println("cap: %d, len: %d", cap(sl1), len(sl1))

    // way #2
    sl2 := []int{10,11,12,13,14}
    sl2 = append(sl2[:elem], sl2[elem+1:]...)
    fmt.Println(sl2)
    fmt.Println("cap: %d, len: %d", cap(sl2), len(sl2))
}

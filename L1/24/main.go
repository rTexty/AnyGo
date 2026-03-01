package main

import (
	"fmt"
	"math"
)


type Point struct {
    x float64
    y float64
}


func NewPoint(x,y float64) Point {
    return Point{
        x,
        y,
    }
}

func (p Point) Distance(anotherPoint Point) float64 {
    dist := math.Sqrt(math.Pow(p.x - anotherPoint.x, 2) + math.Pow(p.y - anotherPoint.y,2))
    return dist
}

func main() {
    p1 := NewPoint(3,4)
    p2 := NewPoint(6,8)
    dist := p1.Distance(p2)
    fmt.Println("Distance is : ", dist)
}
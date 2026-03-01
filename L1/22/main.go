package main

import (
	"fmt"
	"math/big"
)

func main() {
	var int1 big.Int
	var int2 big.Int
	var int3 big.Int
	var int4 big.Int
	var int5 big.Int
	var int6 big.Int
 	fmt.Scan(&int1, &int2)
	int3 = *int3.Add(&int1, &int2)
	int4 = *int4.Mul(&int1,&int2)
	int5 = *int5.Div(&int1,&int2)
	int6 = *int6.Sub(&int1, &int2)
	fmt.Println(&int3)
	fmt.Println(int4)
	fmt.Println(int5)
	fmt.Println(int6)
}
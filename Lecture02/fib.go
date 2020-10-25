package main

import (
	"fmt"
	"math/big"
)

func fib(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func fastFib(n int, memo map[int]big.Int) big.Int {
	if n == 0 || n == 1 {
		return *big.NewInt(1)
	}
	val, ok := memo[n]
	if !ok {
		n1 := fastFib(n-1, memo)
		n2 := fastFib(n-2, memo)
		memo[n] = *(new(big.Int).Add(&n1, &n2))
		return memo[n]
	}
	return val
}

func main() {
	for i := 1; i <= 120; i++ {
		//fmt.Println("fib(", i, ") =", fib(i))
		memo := make(map[int]big.Int)
		var f big.Int
		f = fastFib(i, memo)
		fmt.Printf("fib(%d) = %s\n", i, f.Text(10))
	}
}

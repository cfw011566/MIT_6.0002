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

func fastFib(n *big.Int, memo map[string]big.Int) big.Int {
	if n.Cmp(big.NewInt(0)) == 0 || n.Cmp(big.NewInt(1)) == 0 {
		return *big.NewInt(1)
	}
	val, ok := memo[n.Text(10)]
	if !ok {
		var n1 big.Int
		var n2 big.Int
		n1.Sub(n, big.NewInt(1))
		n2.Sub(n, big.NewInt(2))
		n1 = fastFib(&n1, memo)
		n2 = fastFib(&n2, memo)
		var result big.Int
		result.Add(&n1, &n2)
		memo[n.Text(10)] = result
		return result
	}
	return val
}

func main() {
	for i := 1; i <= 120; i++ {
		//fmt.Println("fib(", i, ") =", fib(i))
		memo := make(map[string]big.Int)
		var f big.Int
		f = fastFib(big.NewInt(int64(i)), memo)
		fmt.Println("fib(", i, ") =", f.Text(10))
	}
}

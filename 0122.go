package main

import (
	"fmt"
)

func main() {
	l := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println(smap(square(), l), smap(cube(), l), smap(sum(), l), smap(average(), l))
}
func smap(f func(int) int, l []int) []int {
	re := make([]int, len(l))
	for i, v := range l {
		re[i] = f(v)
	}
	return re
}
func square() func(x int) int {
	re := func(x int) int {
		return x * x
	}
	return re
}
func cube() func(x int) int {
	re := func(x int) int {
		return x * x * x
	}
	return re
}
func sum() func(x int) int {
	i := 0
	a := func(x int) int {
		i = i + x
		return i
	}
	return a
}
func average() func(x int) int {
	i := 0
	n := 0
	a := func(x int) int {
		i = i + x
		n = n + 1
		//fi := float32(i)
		//fn := float32(n)
		return i / n
	}
	return a
}

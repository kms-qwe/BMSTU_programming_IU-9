/*
Реализация длинной арифметики по заданному основанию

*/

package main

import "fmt"

func main() {
	fmt.Println(add([]int32{1, 1}, []int32{1, 0, 1}, 2))
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func add(a, b []int32, p int) []int32 {
	p32 := int32(p)
	result := make([]int32, max(len(a), len(b)))
	var quotient int32
	for i := 0; i < len(a) || i < len(b); i++ {
		if i >= len(a) {
			sum := b[i] + quotient
			result[i], quotient = sum%p32, sum/p32
		} else if i >= len(b) {
			sum := a[i] + quotient
			result[i], quotient = sum%p32, sum/p32
		} else {
			sum := a[i] + b[i] + quotient
			result[i], quotient = sum%p32, sum/p32
		}
	}
	if quotient > 0 {
		result = append(result, quotient)
	}
	return result
}

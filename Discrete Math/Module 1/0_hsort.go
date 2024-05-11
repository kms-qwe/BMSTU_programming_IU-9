package main

import "fmt"

func main() {
	var arr = []int{
		-158, -173, -194, -175, -277, -267, -196, -189, -275, -204,
		-250, -161, -163, -227, -196, -301, -210, -248, -218, -292,
		-260, -179, -232, -269, -284, -239,
	}
	fmt.Println(arr)
	hsort(len(arr),
		func(i, j int) bool {
			return arr[i] < arr[j]
		},
		func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
	fmt.Println(arr)
}
func hsort(n int,
	less func(i, j int) bool,
	swap func(i, j int)) {

	heapify := func(i, n int) {
		for {
			l, r, j := 2*i+1, 2*i+2, i
			if l < n && less(i, l) {
				i = l
			}
			if r < n && less(i, r) {
				i = r
			}
			if i == j {
				return
			}
			swap(i, j)
		}
	}

	for i := n/2 - 1; i >= 0; i-- {
		heapify(i, n)
	}

	for i := n - 1; i > 0; i-- {
		swap(0, i)
		heapify(0, i)
	}
}

package main

func MergeSort(items int, compare func(i, j int) int, indices chan int) {
	merge := func(left, right chan int, result chan int) {
		leftVal, leftOpen := <-left
		rightVal, rightOpen := <-right

		for leftOpen || rightOpen {
			if !leftOpen {
				result <- rightVal
				rightVal, rightOpen = <-right
			} else if !rightOpen {
				result <- leftVal
				leftVal, leftOpen = <-left
			} else if compare(leftVal, rightVal) <= 0 {
				result <- leftVal
				leftVal, leftOpen = <-left
			} else {
				result <- rightVal
				rightVal, rightOpen = <-right
			}
		}
		close(result)
	}

	var sort func(start, end int, indices chan int)
	sort = func(start, end int, indices chan int) {
		if end-start <= 1 {
			if end-start == 1 {
				indices <- start
			}
			close(indices)
			return
		}

		mid := (start + end) / 2

		leftChan := make(chan int)
		rightChan := make(chan int)

		go sort(start, mid, leftChan)
		go sort(mid, end, rightChan)

		merge(leftChan, rightChan, indices)
	}

	sort(0, items, indices)
}

func main() {

}

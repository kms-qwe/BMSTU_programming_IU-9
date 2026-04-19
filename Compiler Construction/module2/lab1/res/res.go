package main

import "fmt"

func grade(score int) string {
	switch {
	case score ==
		5:
		return "excellent"
	case score ==

		4:
		return "good"
	case score ==

		3:
		return "ok"
	default:
		return "bad"
	}
}

func dayType(day string) string {
	switch {
	case day ==
		"sat" ||
		day ==
			"sun":
		return "weekend"
	case day ==

		"mon" ||
		day ==

			"tue" ||
		day ==

			"wed" ||
		day ==

			"thu" ||
		day ==

			"fri":
		return "workday"
	default:
		return "unknown"
	}
}

func alreadyExprless(x int) string {
	switch {
	case x < 0:
		return "negative"
	case x == 0:
		return "zero"
	default:
		return "positive"
	}
}

func withInit(x int) string {
	switch y := x % 3; {
	case y ==
		0:
		return "divisible by 3"
	case y ==

		1 ||
		y ==

			2:
		return "not divisible by 3"
	default:
		return "impossible"
	}
}

func main() {
	fmt.Println("grade(5):", grade(5))
	fmt.Println("grade(2):", grade(2))

	fmt.Println("dayType(\"sat\"):", dayType("sat"))
	fmt.Println("dayType(\"wed\"):", dayType("wed"))
	fmt.Println("dayType(\"???\"):", dayType("???"))

	fmt.Println("alreadyExprless(-1):", alreadyExprless(-1))
	fmt.Println("alreadyExprless(0):", alreadyExprless(0))
	fmt.Println("alreadyExprless(7):", alreadyExprless(7))

	fmt.Println("withInit(6):", withInit(6))
	fmt.Println("withInit(7):", withInit(7))
}

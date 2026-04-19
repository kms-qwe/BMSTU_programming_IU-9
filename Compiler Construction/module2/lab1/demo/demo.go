package main

import "fmt"

func grade(score int) string {
	switch score {
	case 5:
		return "excellent"
	case 4:
		return "good"
	case 3:
		return "ok"
	default:
		return "bad"
	}
}

func dayType(day string) string {
	switch day {
	case "sat", "sun":
		return "weekend"
	case "mon", "tue", "wed", "thu", "fri":
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
	switch y := x % 3; y {
	case 0:
		return "divisible by 3"
	case 1, 2:
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

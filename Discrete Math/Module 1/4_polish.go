package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	expr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(polishEvaluate(polishWithoutBrackets(expr)))
}
func polishWithoutBrackets(expr string) string {
	expr = strings.ReplaceAll(expr, "(", "")
	expr = strings.ReplaceAll(expr, ")", "")
	expr = strings.ReplaceAll(expr, " ", "")
	expr = strings.ReplaceAll(expr, "\n", "")
	return expr
}
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func isOperand(symb string) bool {
	return symb == "+" || symb == "-" || symb == "*"
}
func polishEvaluate(expr string) int {
	// fmt.Printf("|%s|\n", expr)
	stack := []int{}
	for i := len(expr) - 1; i >= 0; i-- {
		// fmt.Printf("%d %v %v\n", i, string(expr[i]), stack)
		switch s := string(expr[i]); s {
		case "+":
			stack[len(stack)-2] = stack[len(stack)-1] + stack[len(stack)-2]
			stack = stack[:len(stack)-1]
		case "-":
			stack[len(stack)-2] = stack[len(stack)-1] - stack[len(stack)-2]
			stack = stack[:len(stack)-1]
		case "*":
			stack[len(stack)-2] = stack[len(stack)-1] * stack[len(stack)-2]
			stack = stack[:len(stack)-1]
		default:
			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Printf("%s\n", "convert error")
				return -1
			}
			stack = append(stack, num)
		}
	}
	return stack[0]
}

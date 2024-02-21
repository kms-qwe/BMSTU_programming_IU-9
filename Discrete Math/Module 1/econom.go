package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	expr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(countOperation(polishWithoutBrackets(expr)))
}
func polishWithoutBrackets(expr string) string {
	expr = strings.ReplaceAll(expr, "(", "")
	expr = strings.ReplaceAll(expr, ")", "")
	expr = strings.ReplaceAll(expr, " ", "")
	expr = strings.ReplaceAll(expr, "\n", "")
	return expr
}

func isOperand(symb string) bool {
	return symb == "#" || symb == "$" || symb == "@"
}
func countOperation(expr string) int {
	stack := []string{}
	alreadyCalculate := map[string]bool{}
	count := 0
	for i := len(expr) - 1; i >= 0; i-- {
		if s := string(expr[i]); "a" <= s && s <= "z" {
			stack = append(stack, s)
		} else {
			stack[len(stack)-2] = s + stack[len(stack)-1] + stack[len(stack)-2]
			stack = stack[:len(stack)-1]
			if _, ok := alreadyCalculate[stack[len(stack)-1]]; !ok {
				count += 1
				alreadyCalculate[stack[len(stack)-1]] = true
			}
		}
	}
	return count

}

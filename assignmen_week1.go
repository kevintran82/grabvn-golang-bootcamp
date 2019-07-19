package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parse(input string) (a float64, b float64, op string, err error) {
	params := strings.Fields(input)

	if len(params) < 3 {
		return 0.0, 0.0, "", errors.New("Error: Invalid input")
	}

	x, err1 := strconv.ParseFloat(params[0], 64)
	if err1 != nil {
		return 0.0, 0.0, "", err1
	}

	y, err2 := strconv.ParseFloat(params[2], 64)
	if err2 != nil {
		return 0.0, 0.0, "", err2

	}

	operator := params[1]

	return x, y, operator, nil
}

func calculate(input string) (result float64, err error) {
	total := 0.0
	a, b, op, err := parse(input)
	if err != nil {
		return total, err
	}

	switch op {
	case "+":
		total = a + b
	case "-":
		total = a - b
	case "*":
		total = a * b
	case "/":
		if b == 0 {
			return total, errors.New("Error: Can't divide by zero")
		}

		total = a / b
	default:
		return 0, errors.New("Error: Invalid operator")
	}

	return total, nil
}

func main() {

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for reader.Scan() {
		input := reader.Text()

		result, err := calculate(input)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(input, " = ", result)
		}

		fmt.Print("> ")
	}
}

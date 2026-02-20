package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() (T, error) {
	stackLength := len(s.values)
	if stackLength == 0 {
		var null T
		return null, errors.New("Can't pop from an empty Stack")
	}

	res := s.values[stackLength-1]
	s.values = s.values[:stackLength-1]

	return res, nil
}

func (s *Stack[T]) Peek() (T, error) {
	stackLength := len(s.values)
	if stackLength == 0 {
		var null T
		return null, errors.New("Can't peek on an empty Stack")
	}

	return s.values[stackLength-1], nil
}

func (s *Stack[T]) Len() int {
	return len(s.values)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("No equation were passed")
		return
	}

	equation := args[0]

	result, err := evaluate(equation)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}

func evaluate(equation string) (float64, error) {
	values := Stack[float64]{}
	operators := Stack[rune]{}

	for i := 0; i < len(equation); i++ {
		character := rune(equation[i])

		switch {
		case isDigit(character):
			start := i
			for i < len(equation) && (isDigit(rune(equation[i])) || rune(equation[i]) == '.') {
				i++
			}

			numberString := equation[start:i]
			val, err := strconv.ParseFloat(numberString, 64)

			if err != nil {
				return 0, err
			}
			values.Push(val)

			i--

		case character == '(':
			operators.Push(character)

		case character == ')':
			operator, err := operators.Peek()
			for operator != '(' && err == nil {
				err := executeTopOperator(&values, &operators)
				if err != nil {
					return 0, err
				}
			}

			if err != nil {
				return 0, err
			} else {
				operators.Pop()
			}

		case isOperator(character):
			operator, err := operators.Peek()
			for priority(operator) >= priority(character) && err == nil {
				err := executeTopOperator(&values, &operators)
				if err != nil {
					return 0, err
				}
			}

			if err != nil {
				return 0, err
			}

			operators.Push(character)
		}
	}

	for operators.Len() > 0 {
		err := executeTopOperator(&values, &operators)
		if err != nil {
			return 0, err
		}
	}

	result, err := values.Pop()

	if err != nil {
		return 0, errors.New("Incomplete expression")
	}

	if values.Len() > 0 {
		return 0, errors.New("Too many values in the expression")
	}

	return result, nil
}

func isDigit(character rune) bool {
	return character >= '0' && character <= '9'
}

func isOperator(character rune) bool {
	return character == '+' || character == '-' || character == '*' || character == '/'
}

func priority(operator rune) int {
	switch operator {
	case '+', '-':
		return 1

	case '*', '/':
		return 2
	}

	return 0
}

func executeTopOperator(values *Stack[float64], operators *Stack[rune]) error {
	if operators.Len() < 1 {
		return errors.New("No operators to execute")
	}
	operator, _ := operators.Pop()

	if values.Len() < 2 {
		return fmt.Errorf("Not enough values to use with '%c'", operator)
	}
	secondValue, _ := values.Pop()
	firstValue, _ := values.Pop()

	var newValue float64
	switch operator {
	case '+':
		newValue = firstValue + secondValue
	case '-':
		newValue = firstValue - secondValue
	case '*':
		newValue = firstValue * secondValue
	case '/':
		if secondValue == 0 {
			return errors.New("Division by zero")
		}

		newValue = firstValue / secondValue
	}

	values.Push(newValue)

	return nil
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

const groupOpen = '('
const groupClose = ')'

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() (T, error) {
	stackLength := len(s.values)
	if stackLength == 0 {
		var zero T
		return zero, errors.New("Can't pop from an empty Stack")
	}

	res := s.values[stackLength-1]
	s.values = s.values[:stackLength-1]

	return res, nil
}

func (s *Stack[T]) Peek() (T, error) {
	stackLength := len(s.values)
	if stackLength == 0 {
		var zero T
		return zero, errors.New("Can't peek on an empty Stack")
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
		return
	}

	fmt.Println(result)
}

func evaluate(expression string) (float64, error) {
	const ignoreCharacter = ' '
	const decimalSeparator = '.'

	values := Stack[float64]{}
	operators := Stack[rune]{}

	for i := 0; i < len(expression); i++ {
		character := rune(expression[i])

		switch {
		case character == ignoreCharacter:
			continue

		case isDigit(character):
			start := i
			for i < len(expression) && (isDigit(rune(expression[i])) || rune(expression[i]) == decimalSeparator) {
				i++
			}

			numberString := expression[start:i]
			val, err := strconv.ParseFloat(numberString, 64)

			if err != nil {
				return 0, err
			}

			values.Push(val)

			i--

		case character == groupOpen:
			operators.Push(character)

		case character == groupClose:
			operator, err := operators.Peek()

			if err != nil {
				return 0, fmt.Errorf("Mismatched parentheses: unexpected '%c'", groupClose)
			}

			if expression[i-1] == groupOpen {
				return 0, errors.New("Empty parentheses")
			}

			for operator != groupOpen {
				err = executeTopOperator(&values, &operators)
				if err != nil {
					return 0, err
				}

				operator, err = operators.Peek()
			}

			operators.Pop()

		case isOperator(character):
			operator, err := operators.Peek()

			for operators.Len() > 0 && priority(operator) >= priority(character) {
				err = executeTopOperator(&values, &operators)
				if err != nil {
					return 0, err
				}

				operator, err = operators.Peek()
			}

			operators.Push(character)

		default:
			return 0, fmt.Errorf("Unexpected character '%c'", character)
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

func ParseFloat64(initialString string) (float64, error) {
	return strconv.ParseFloat(initialString, 64)
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

	if operator == groupOpen {
		return fmt.Errorf("Mismatched parentheses: unclosed '%c'", groupOpen)
	}

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

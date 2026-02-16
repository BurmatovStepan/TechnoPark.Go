package main

import (
	"fmt"
	"flag"
)

type stack[T any] struct {
	values []T
}


func (s *stack[T]) Push(value T) {
	s.values = append(s.values, value)
}


func (s *stack[T]) Pop() T {
	l := len(s.values)
    res := s.values[l - 1]
    s.values = s.values[:l - 1]

    return res
}


func (s *stack[T]) Len() int {
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

	result := evaluate(equation)
	fmt.Println(result)
}


func evaluate(equation string) int {
	values :=
	return 0
}

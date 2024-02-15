package stack

import "errors"

type (
	Stack[T any] struct {
		top    *node[T]
		length int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

var (
	ErrEmpty = errors.New("empty")
)

// Create a new stack
func New[T any]() *Stack[T] {
	return &Stack[T]{nil, 0}
}

// Return the number of items in the stack
func (stack *Stack[T]) Len() int {
	return stack.length
}

// View the top item on the stack
func (stack *Stack[T]) Peek() (T, error) {
	if stack.length == 0 {
		return *new(T), ErrEmpty
	}
	return stack.top.value, nil
}

// Pop the top item of the stack and return it
func (stack *Stack[T]) Pop() (T, error) {
	if stack.length == 0 {
		return *new(T), ErrEmpty
	}

	n := stack.top
	stack.top = n.prev
	stack.length--
	return n.value, nil
}

// Push a value onto the top of the stack
func (stack *Stack[T]) Push(value T) {
	n := &node[T]{value, stack.top}
	stack.top = n
	stack.length++
}

// Push a value onto the top of the stack
func (stack *Stack[T]) ToSlice() []T {
	result := make([]T, 0, stack.length)

	pointer := stack.top
	for pointer != nil {
		result = append(result, pointer.value)
		pointer = pointer.prev
	}

	return result
}

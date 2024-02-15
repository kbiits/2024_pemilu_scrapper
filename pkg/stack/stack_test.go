package stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	stack := New[string]()
	stack.Push("abc")
	stack.Push("def")
	stack.Push("ghi")
	stack.Push("jkl")

	stackInSlice := stack.ToSlice()
	require.EqualValues(t, []string{"jkl", "ghi", "def", "abc"}, stackInSlice)

	stack.Pop()
	stack.Pop()
	stack.Pop()
	stackInSlice = stack.ToSlice()
	require.EqualValues(t, 1, stack.length)
	require.EqualValues(t, []string{"abc"}, stackInSlice)

	peekVal, err := stack.Peek()
	require.NoError(t, err)
	require.EqualValues(t, "abc", peekVal)
	require.Equal(t, 1, stack.Len())
}

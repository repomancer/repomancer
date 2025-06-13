package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobQueue_Add_and_Pop(t *testing.T) {
	q := NewJobQueue()
	q.Add(&Job{Command: "1"})
	q.Add(&Job{Command: "2"})
	assert.Equal(t, 2, len(q.items))
	j := q.Pop()
	assert.Equal(t, &Job{Command: "1"}, j)
	j = q.Pop()
	assert.Equal(t, &Job{Command: "2"}, j)
	j = q.Pop()
	assert.Nil(t, j)
}

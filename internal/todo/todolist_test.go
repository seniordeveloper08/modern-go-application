package todo

import (
	"context"
	"github.com/goph/idgen"
	"github.com/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type todoEventsStub struct {
	markedAsDone MarkedAsDone
}

func (s *todoEventsStub) MarkedAsDone(ctx context.Context, event MarkedAsDone) error {
	s.markedAsDone = event

	return nil
}

func TestTodoList_CreatesATodo(t *testing.T) {
	todoStore := NewInmemoryTodoStore()
	todoList := NewTodoList(idgen.NewConstantGenerator("id"), todoStore, nil)

	req := CreateTodoRequest{
		Text: "My first todo",
	}

	resp, err := todoList.CreateTodo(context.Background(), req)
	require.NoError(t, err)

	expectedResponse := &CreateTodoResponse{
		ID: "id",
	}

	assert.Equal(t, expectedResponse, resp)

	expectedTodo := Todo{
		ID:   resp.ID,
		Text: req.Text,
	}

	todo, err := todoStore.Get(context.Background(), resp.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, todo)
}

func TestTodoList_CannotCreateATodo(t *testing.T) {
	todoList := NewTodoList(idgen.NewConstantGenerator("id"), NewReadOnlyTodoStore(NewInmemoryTodoStore()), nil)

	req := CreateTodoRequest{
		Text: "My first todo",
	}

	_, err := todoList.CreateTodo(context.Background(), req)
	require.Error(t, err)
}

func TestTodoList_ListTodos(t *testing.T) {
	todoStore := NewInmemoryTodoStore()

	todo := Todo{
		ID:   "id",
		Text: "Make the listing work",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	todoList := NewTodoList(idgen.NewConstantGenerator("id"), todoStore, nil)

	todos, err := todoList.ListTodos(context.Background())
	require.NoError(t, err)

	expectedTodos := &ListTodosResponse{
		Todos: []Todo{
			todo,
		},
	}

	assert.Equal(t, expectedTodos, todos)
}

func TestTodoList_MarkAsDone(t *testing.T) {
	todoStore := NewInmemoryTodoStore()

	todo := Todo{
		ID:   "id",
		Text: "Do me",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	events := &todoEventsStub{}
	todoList := NewTodoList(nil, todoStore, events)

	req := MarkAsDoneRequest{
		ID: "id",
	}

	err := todoList.MarkAsDone(context.Background(), req)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Done = true

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)

	expectedEvent := MarkedAsDone{
		ID: "id",
	}

	assert.Equal(t, expectedEvent, events.markedAsDone)
}

func TestTodoList_CannotMarkANonExistingTodoDone(t *testing.T) {
	todoStore := NewInmemoryTodoStore()

	events := &todoEventsStub{}
	todoList := NewTodoList(nil, todoStore, events)

	req := MarkAsDoneRequest{
		ID: "id",
	}

	err := todoList.MarkAsDone(context.Background(), req)
	require.Error(t, err)

	cause := errors.Cause(err)

	require.IsType(t, TodoNotFoundError{}, cause)

	e := cause.(TodoNotFoundError)
	assert.Equal(t, "id", e.ID)
}

func TestTodoList_StoringDoneTodoFails(t *testing.T) {
	inmemTodoStore := NewInmemoryTodoStore()

	todo := Todo{
		ID:   "id",
		Text: "Do me",
	}
	require.NoError(t, inmemTodoStore.Store(context.Background(), todo))

	todoList := NewTodoList(nil, NewReadOnlyTodoStore(inmemTodoStore), &todoEventsStub{})

	req := MarkAsDoneRequest{
		ID: "id",
	}

	err := todoList.MarkAsDone(context.Background(), req)
	require.Error(t, err)
}

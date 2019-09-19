package todo

import (
	"context"

	"emperror.dev/errors"
)

// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// List manages a list of todos.
type List struct {
	idgenerator IDGenerator
	store       Store
	events      Events
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// Store stores todos.
type Store interface {
	// Store stores a todo.
	Store(ctx context.Context, todo Todo) error

	// All returns all todos.
	All(ctx context.Context) ([]Todo, error)

	// Get returns a single todo by its ID.
	Get(ctx context.Context, id string) (Todo, error)
}

// NotFoundError is returned if a todo cannot be found.
type NotFoundError struct {
	ID string
}

// Error implements the error interface.
func (NotFoundError) Error() string {
	return "todo not found"
}

// Details returns error details.
func (e NotFoundError) Details() []interface{} {
	return []interface{}{"todo_id", e.ID}
}

// IsBusinessError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (NotFoundError) IsBusinessError() bool {
	return true
}

//go:generate sh -c "test -x ${MGA} && ${MGA} gen ev dispatcher --from Events"

// Events dispatches todo events.
type Events interface {
	// MarkedAsDone dispatches a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event MarkedAsDone) error
}

// MarkedAsDone event is triggered when a todo gets marked as done.
type MarkedAsDone struct {
	ID string
}

// NewList returns a new todo list.
func NewList(id IDGenerator, todos Store, events Events) *List {
	return &List{
		idgenerator: id,
		store:       todos,
		events:      events,
	}
}

// CreateTodo adds a new todo to the list.
func (l *List) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := l.idgenerator.Generate()
	if err != nil {
		return "", err
	}

	todo := Todo{
		ID:   id,
		Text: text,
	}

	err = l.store.Store(ctx, todo)

	return id, err
}

// ListTodos returns the list of todos.
func (l *List) ListTodos(ctx context.Context) ([]Todo, error) {
	return l.store.All(ctx)
}

// MarkAsDone marks a todo as done.
func (l *List) MarkAsDone(ctx context.Context, id string) error {
	todo, err := l.store.Get(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	todo.Done = true

	err = l.store.Store(ctx, todo)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	event := MarkedAsDone{
		ID: todo.ID,
	}

	err = l.events.MarkedAsDone(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	return nil
}

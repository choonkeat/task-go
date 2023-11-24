// Task is a simple wrapper around a data and an error
// It is inspired by typed functional languages' Result
// to make error handling composable.
package task

// Task is a simple wrapper around a data and an error
type Task[T interface{}] struct {
	data T
	err  error
}

// Ok function creates a task with data and no error
func Ok[T interface{}](data T) *Task[T] {
	return &Task[T]{data, nil}
}

// Err function creates a task with error and no data
func Err[T interface{}](err error) *Task[T] {
	return &Task[T]{err: err}
}

// Wrap function creates a task with data and error
func Wrap[T interface{}](data T, err error) *Task[T] {
	return &Task[T]{data, err}
}

// Scenarios is a struct that holds the callbacks for the success and failure scenarios
type Scenarios[T interface{}] struct {
	Ok  func(data T)
	Err func(err error)
}

// Unwrap method allows us to unambiguously handle the mutually exclusive
// scenarios of success and failure
func (t *Task[T]) Unwrap(s Scenarios[T]) {
	if t.err != nil {
		s.Err(t.err)
	} else {
		s.Ok(t.data)
	}
}

// AndThen method returns an error if the task is already in error state
// Otherwise, pass the data to the callback to change our _task_
func (t *Task[T]) AndThen(callback func(T) *Task[T]) *Task[T] {
	if t.err != nil {
		return t
	}
	return callback(t.data)
}

// AndThen function is like AndThen _method_ but lets you change the type of the data
func AndThen[T, Y any](t *Task[T], callback func(T) *Task[Y]) *Task[Y] {
	if t.err != nil {
		return &Task[Y]{err: t.err}
	}
	return callback(t.data)
}

// Map method returns an error if the task is already in error state
// Otherwise, pass the data to the callback to change our _data_
func (t *Task[T]) Map(callback func(T) T) *Task[T] {
	if t.err != nil {
		return t
	}
	return &Task[T]{callback(t.data), nil}
}

// Map function is like Map _method_ but lets you change the type of the data
func Map[T, Y any](t *Task[T], callback func(T) Y) *Task[Y] {
	if t.err != nil {
		return &Task[Y]{err: t.err}
	}
	return &Task[Y]{data: callback(t.data), err: nil}
}

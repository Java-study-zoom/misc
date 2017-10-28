package psqlx

// Maker a make function that creates an interface{} for unmarshaling.
type Maker func() interface{}

// Doer is a functions that performs the a particular task.
type Doer func(v interface{}) error

// Iter is an interator.
type Iter struct {
	Make Maker
	Do   Doer
}

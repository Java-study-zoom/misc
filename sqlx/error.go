package sqlx

import (
	"fmt"
)

type queryError struct {
	q   string
	err error
}

func (e *queryError) Error() string {
	return fmt.Sprintf("%s: query:\n%q", e.err, e.q)
}

// Error creates a query error if err is not nil.
// It returns nil if err is nil
func Error(q string, err error) error {
	if err == nil {
		return nil
	}
	return &queryError{q, err}
}

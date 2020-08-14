package errors

type Errors []error

func (ee Errors) Error() string {
	s := ""

	for _, err := range ee {
		s += err.Error() + ";"
	}

	return s
}

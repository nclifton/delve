package valid

import "strings"

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Errors returns itself.
func (es Errors) Errors() []error {
	return es
}

func (es Errors) Error() string {
	var errs []string
	for _, e := range es {
		errs = append(errs, e.Error())
	}
	return strings.Join(errs, "; ")
}

// Error encapsulates a name, an error and whether there's a custom error message or not.
type Error struct {
	Name string
	Err  error

	// Validator indicates the name of the validator that failed
	Validator string
	Path      []string
}

func (e Error) Error() string {
	errName := e.Name
	if len(e.Path) > 0 {
		errName = strings.Join(append(e.Path, e.Name), ".")
	}

	return errName + ": " + e.Err.Error()
}

func PrependPathToErrors(err error, path string) error {
	switch err2 := err.(type) {
	case Error:
		err2.Path = append([]string{path}, err2.Path...)
		return err2
	case Errors:
		errors := err2.Errors()
		for i, err3 := range errors {
			errors[i] = PrependPathToErrors(err3, path)
		}
		return err2
	}
	return err
}

// ErrorByField returns error for specified field of the struct
// validated by ValidateStruct or empty string if there are no errors
// or this field doesn't exists or doesn't have any errors.
func ErrorByField(e error, field string) string {
	if e == nil {
		return ""
	}
	return ErrorsByField(e)[field]
}

// ErrorsByField returns map of errors of the struct validated
// by ValidateStruct or empty map if there are no errors.
func ErrorsByField(e error) map[string]string {
	m := make(map[string]string)
	if e == nil {
		return m
	}
	// prototype for ValidateStruct

	switch e := e.(type) {
	case Error:
		m[e.Name] = e.Err.Error()
	case Errors:
		for _, item := range e.Errors() {
			n := ErrorsByField(item)
			for k, v := range n {
				m[k] = v
			}
		}
	}

	return m
}

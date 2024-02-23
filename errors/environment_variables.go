package errors

import "fmt"

type EnvironmentVariableNotFoundError struct {
	EnvironmentVariable string
}

func NewEnvironmentVariableNotFoundError(environmentVariable string) *EnvironmentVariableNotFoundError {
	return &EnvironmentVariableNotFoundError{
		EnvironmentVariable: environmentVariable,
	}
}

func (e *EnvironmentVariableNotFoundError) Error() string {
	return fmt.Sprintf("Environment variable not found: %s", e.EnvironmentVariable)
}

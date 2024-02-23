package errors

import "fmt"

type UnknownDatabaseTypeError struct {
	Err string
}

func NewUnknownDatabaseTypeError(err string) *UnknownDatabaseTypeError {
	return &UnknownDatabaseTypeError{
		Err: err,
	}
}

func (e *UnknownDatabaseTypeError) Error() string {
	return fmt.Sprintf("Unknown database type: %s", e.Err)
}

type DatabaseConnectionError struct {
	Err string
}

func NewDatabaseConnectionError(err string) *DatabaseConnectionError {
	return &DatabaseConnectionError{
		Err: err,
	}
}

func (e *DatabaseConnectionError) Error() string {
	return fmt.Sprintf("Database connection error: %s", e.Err)
}

type DatabaseMigrationError struct {
	Err string
}

func NewDatabaseMigrationError(err string) *DatabaseMigrationError {
	return &DatabaseMigrationError{
		Err: err,
	}
}

func (e *DatabaseMigrationError) Error() string {
	return fmt.Sprintf("Database migration error: %s", e.Err)
}

type ObjectNotFoundError struct {
	Object string
}

func NewObjectNotFoundError(object string) *ObjectNotFoundError {
	return &ObjectNotFoundError{
		Object: object,
	}
}

func (e *ObjectNotFoundError) Error() string {
	return fmt.Sprintf("Object not found: %s", e.Object)
}

type MainAdminDeletionError struct {
	Err string
}

func NewMainAdminDeletionError(err string) *MainAdminDeletionError {
	return &MainAdminDeletionError{
		Err: err,
	}
}

func (e *MainAdminDeletionError) Error() string {
	return fmt.Sprintf("Main admin deletion error: %s", e.Err)
}

type SlugNotUniqueError struct {
	Slug string
}

func NewSlugNotUniqueError(slug string) *SlugNotUniqueError {
	return &SlugNotUniqueError{
		Slug: slug,
	}
}

func (e *SlugNotUniqueError) Error() string {
	return fmt.Sprintf("Slug not unique: %s", e.Slug)
}

package project

// NotFoundError represents account not found error
type NotFoundError struct{}

func (anf NotFoundError) Error() string {
	return "Project Not Found"
}

// AlreadyExistError represents account already exist error
type AlreadyExistError struct{}

func (aee AlreadyExistError) Error() string {
	return "Project already exist"
}

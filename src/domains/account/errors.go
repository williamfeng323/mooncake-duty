package account

// NotFoundError represents account not found error
type NotFoundError struct{}

func (anf NotFoundError) Error() string {
	return "Account Not Found"
}

// AlreadyExistError represents account already exist error
type AlreadyExistError struct{}

func (aee AlreadyExistError) Error() string {
	return "Account already exist"
}

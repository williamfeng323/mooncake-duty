package issue

// NotFoundError represents account not found error
type NotFoundError struct{}

func (inf NotFoundError) Error() string {
	return "Issue Not Found"
}

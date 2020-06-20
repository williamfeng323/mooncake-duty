package shift

// NoMemberError represents account already exist error
type NoMemberError struct{}

func (n NoMemberError) Error() string {
	return "No members in the shift"
}

// OutOfScopeError represents period out of scope error
type OutOfScopeError struct{}

func (o OutOfScopeError) Error() string {
	return "Selected period out of the shift scope"
}

// DuplicateShiftError represents shift can only associate with project that have no shift configured
type DuplicateShiftError struct{}

func (d DuplicateShiftError) Error() string {
	return "Shift can only associate with project that have no shift configured"
}

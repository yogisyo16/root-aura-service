package services

// Models is a convenience aggregate of the domain services. Not currently
// constructed or used anywhere — api/main.go wires Todo, User, and
// TodoDetails into the handlers individually instead. Kept here as a
// reference for "these are the domain models this app has."
type Models struct {
	Todo    Todo
	User    User
	Details TodoDetails
}

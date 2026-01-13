package command

type UpdateUserCommand struct {
	UserID    string
	Email     *string
	FirstName *string
	LastName  *string
}
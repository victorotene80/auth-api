package response

type UserResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalCount int64          `json:"total_count"`
}

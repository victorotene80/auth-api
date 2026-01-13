package query

type ListUsersQuery struct {
	Page     int
	PageSize int
	Role     *string 
	IsActive *bool   
}

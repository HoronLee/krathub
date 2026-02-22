package entity

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
	Phone    *string
	Avatar   *string
	Bio      *string
	Location *string
	Website  *string
	Role     string
}

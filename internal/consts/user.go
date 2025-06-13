package consts

type UserRole uint

const (
	Guest UserRole = iota + 1
	User
	Admin
	Operator
)

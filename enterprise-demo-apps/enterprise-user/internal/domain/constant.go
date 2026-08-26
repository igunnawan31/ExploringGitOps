package domain

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleCashier  UserRole = "cashier"
	RoleCustomer UserRole = "customer"
	RoleManager  UserRole = "manager"
)

type UserStatus string

const (
	Active   UserStatus = "active"
	Inactive UserStatus = "inactive"
)

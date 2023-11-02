package models

type UserRole string

const (
	OwnerUserRole UserRole = "OWNER"
	UserUserRole  UserRole = "USER"
)

type UserStatus string

const (
	AvailableUserStatus UserStatus = "AVAILABLE"
	ReservedUserStatus  UserStatus = "RESERVED"
	CheckedInUserStatus UserStatus = "CHECKED-IN"
)

type User struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Name     string     `json:"name"`
	Surname  string     `json:"surname"`
	Role     UserRole   `json:"role"`
	Tel      string     `json:"tel"`
	Status   UserStatus `json:"status" default:"AVAILABLE"`
}

package domain

type Role string

const (
	Roleadmin   Role = "ADMIN"
	RoleService Role = "SERVICE"
)

type User struct {
	ID       string
	EMAIL    string
	PASSHASH string
	ROLE     string
}

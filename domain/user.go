package domain

const (
	RoleAdmin = iota
	RoleUser
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	FullName string `json:"full_name" bson:"full_name"`
	Role     int    `json:"role" bson:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

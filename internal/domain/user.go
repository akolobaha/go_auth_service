package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID       primitive.ObjectID `json:"id"`
	Login    string             `json:"login"`
	Password string             `json:"password"`
	Name     string             `json:"name"`
	Role     string             `json:"role"`
	Active   bool               `json:"active"`
	Age      int                `json:"age"`
	Email    string             `json:"email"`
}

type UserInfo struct {
	ID    primitive.ObjectID `json:"id"`
	Name  string             `json:"name"`
	Age   int                `json:"age"`
	Role  string             `json:"role"`
	Email string             `json:"email"`
}

type UserRole struct {
	ID   primitive.ObjectID `json:"id"`
	Role string             `json:"active"`
}

type UserPassword struct {
	ID       primitive.ObjectID `json:"id"`
	Password string             `json:"password"`
}

type UserResetPassword struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserToken struct {
	UserId primitive.ObjectID `json:"id"`
	Token  string             `json:"token"`
}

type UserIsActive struct {
	ID     primitive.ObjectID `json:"id"`
	Active bool               `json:"active"`
}

type UserMessage struct {
	UserId  primitive.ObjectID `json:"id"`
	Message string             `json:"message"`
}

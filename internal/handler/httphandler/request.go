package httphandler

import "authservice/internal/domain"

type SetUserInfoReq struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type SetUserRoleReq struct {
	Role string `json:"role"`
}

type SetUserActiveReq struct {
	Active bool `json:"active"`
}

type ChangePswReq struct {
	Password string `json:"password"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

func (r SetUserRoleReq) IsValid() bool {
	return r.Role == domain.UserRoleDefault || r.Role == domain.UserRoleAdmin
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}

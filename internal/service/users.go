package service

import (
	"authservice/internal/domain"
	"authservice/internal/repository/tokendb"
	"authservice/internal/repository/userdb"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

var users userdb.DB
var tokens tokendb.DB

func Init(userDB userdb.DB, tokenDB tokendb.DB) {
	users = userDB
	tokens = tokenDB
}

func SignUp(lp *domain.LoginPassword) (*domain.UserToken, error) {

	if _, ok := users.CheckExistLogin(lp.Login); ok {
		return nil, errors.New("login " + lp.Login + " already exists")
	}

	newUser := domain.User{
		ID:       primitive.NewObjectID(),
		Login:    lp.Login,
		Password: hash(lp.Password),
		Role:     domain.UserRoleDefault,
		Email:    lp.Email,
		Active:   true,
	}

	if err := users.SetUser(&newUser); err != nil {
		return nil, err
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, newUser.ID); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: newUser.ID,
		Token:  token,
	}, nil
}

func SignIn(lp *domain.LoginPassword) (*domain.UserToken, error) {

	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if user.Password != hash(lp.Password) {
		return nil, errors.New("wrong password")
	}

	token := createToken(lp.Login)

	if err := tokens.SetUserToken(token, *userId); err != nil {
		return nil, err
	}

	return &domain.UserToken{
		UserId: *userId,
		Token:  token,
	}, nil
}

func ResetPassword(lp *domain.LoginPassword) (*domain.UserMessage, error) {

	userId, ok := users.CheckExistLogin(lp.Login)
	if !ok {
		return nil, errors.New("user not found")
	}

	user, err := users.GetUser(*userId)
	if err != nil {
		return nil, err
	}

	if user.Password != hash(lp.Password) {
		return nil, errors.New("wrong password")
	}
	if user.Email == "" {
		return nil, errors.New("email is empty")
	}

	newRandPassword := hash(user.Password[:rand.Intn(64)])[:6]
	user.Password = hash(newRandPassword)

	sendEmail(user.Email, "Новый пароль", newRandPassword)

	return &domain.UserMessage{
		UserId:  *userId,
		Message: "Новый пароль отправлен на email: " + user.Email,
	}, nil
}

func SetUserInfo(ui *domain.UserInfo) error {

	user, err := users.GetUser(ui.ID)
	if err != nil {
		return err
	}

	user.Name = ui.Name
	user.Age = ui.Age
	user.Email = ui.Email

	return users.SetUser(user)

}

func SetUserRole(ui *domain.UserRole) error {

	user, err := users.GetUser(ui.ID)
	if err != nil {
		return err
	}

	user.Role = ui.Role

	return users.SetUser(user)

}

func SetUserIsActive(ui *domain.UserIsActive) error {

	user, err := users.GetUser(ui.ID)
	if err != nil {
		return err
	}

	user.Active = ui.Active

	return users.SetUser(user)

}

func ChangePsw(up *domain.UserPassword) error {

	user, err := users.GetUser(up.ID)
	if err != nil {
		return err
	}

	user.Password = hash(up.Password)

	return users.SetUser(user)
}

func GetUserShortInfo(id primitive.ObjectID) (*domain.UserInfo, error) {

	user, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	ui := domain.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Age:   user.Age,
		Role:  user.Role,
		Email: user.Email,
	}

	return &ui, nil
}

func GetUserIsActive(id primitive.ObjectID) (*domain.UserIsActive, error) {
	user, err := users.GetUser(id)
	if err != nil {
		return nil, err
	}

	ui := domain.UserIsActive{
		ID:     user.ID,
		Active: user.Active,
	}

	return &ui, nil
}

func GetUserFullInfo(id primitive.ObjectID) (*domain.User, error) {

	user, err := users.GetUser(id)
	return user, err
}

func GetUserIDByToken(token string) (*primitive.ObjectID, error) {
	return tokens.GetUserByToken(token)
}

func hash(str string) string {
	hp := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hp[:])
}

func createToken(login string) string {

	timeChs := md5.Sum([]byte(time.Now().String()))
	loginChs := md5.Sum([]byte(login))

	return hex.EncodeToString(timeChs[:]) + hex.EncodeToString(loginChs[:])
}

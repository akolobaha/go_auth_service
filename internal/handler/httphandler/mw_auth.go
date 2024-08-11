package httphandler

import (
	"authservice/internal/service"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		token := req.Header.Get(HeaderAuthorization)

		if len(token) == 0 {
			resp.WriteHeader(http.StatusUnauthorized)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("token is missing"))
			resp.Write(respBody.Marshall())

			return
		}

		userID, err := service.GetUserIDByToken(token)
		if err != nil {
			resp.WriteHeader(http.StatusUnauthorized)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("wrong token"))
			resp.Write(respBody.Marshall())

			return
		}

		req.Header.Set(HeaderUserID, userID.Hex())

		next.ServeHTTP(resp, req)
	})
}

func IsActive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		id := req.Header.Get(HeaderUserID)
		userID, _ := primitive.ObjectIDFromHex(id)

		userInfo, err := service.GetUserIsActive(userID)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("can not find user by id"))
			resp.Write(respBody.Marshall())

			return
		}

		if userInfo.Active != true {
			resp.WriteHeader(http.StatusForbidden)

			respBody := &HTTPResponse{}
			respBody.SetError(errors.New("user is not active"))
			resp.Write(respBody.Marshall())

			return
		}

		next.ServeHTTP(resp, req)
	})
}

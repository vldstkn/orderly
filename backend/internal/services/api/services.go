package api

import (
	"net/http"
)

type ApiService struct {
	JWTSecret string
}

func NewApiService(JWTSecret string) *ApiService {
	return &ApiService{
		JWTSecret: JWTSecret,
	}
}

func (service *ApiService) AddCookie(w *http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   maxAge,
	}
	http.SetCookie(*w, cookie)
}

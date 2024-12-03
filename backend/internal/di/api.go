package di

import "net/http"

type IApiService interface {
	AddCookie(w *http.ResponseWriter, name, value string, maxAge int)
}

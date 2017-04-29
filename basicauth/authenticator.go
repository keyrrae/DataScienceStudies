package basicauth

import (
	"encoding/base64"
	"github.com/keyrrae/monimenta_backend/controllers"
	"net/http"
	"strings"
)

type BasicAuthenticator struct {
	UserController *controllers.UserController
	Credentials    map[string]string
}

func NewBasicAuthenticator(uc *controllers.UserController, cred map[string]string) *BasicAuthenticator {
	return &BasicAuthenticator{
		UserController: uc,
		Credentials:    cred,
	}
}

func (au *BasicAuthenticator) FuncAuth(uid, pass string) bool {
	if val, ok := au.Credentials[uid]; ok {
		return pass == val
	}

	user, err := au.UserController.GetUserFromDB(uid)
	if err != nil {
		return false
	}

	if user.Cred == pass {
		au.Credentials[uid] = pass
		return true
	}

	return false
}

func (au *BasicAuthenticator) FuncAdminAuth(uid, pass string) bool {
	// TODO: implement admin auth functionality
	return uid == "admin" && pass == "admin"
}

// static method
func WrapAuthenticator(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func (au *BasicAuthenticator) UserAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if !au.FuncAuth(pair[0], pair[1]) {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func (au *BasicAuthenticator) AdminAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if !au.FuncAuth(pair[0], pair[1]) {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/STreeChin/contactapi/internal/service"
	"github.com/sirupsen/logrus"
)

//Auth interface
type Auth interface {
	Middleware(next http.Handler) http.Handler
}

type auth struct {
	log *logrus.Logger
	rep service.Repository
}

//NewAuth new auth
func NewAuth(log *logrus.Logger, rep service.Repository) *auth {
	return &auth{log, rep}
}

func (a *auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Endpoints that don't require authentication
		needAuthPaths := []string{"/v1/user/login"}
		for _, value := range needAuthPaths {
			if value == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		apiKey := r.Header.Get("autopilotapikey")
		if apiKey == "" {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request", "message": "No autopilotapikey header provided."})
			if err != nil {
				a.log.Error("Middleware: json encode error")
			}
			return
		}

		_, err := a.rep.GetContactIDByAPIKey(apiKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			err = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized", "message": "Provided autopilotapikey not valid."})
			if err != nil {
				a.log.Error("Middleware: json encode error")
			}
			return
		}

		//ctx := context.WithValue(r.Context(), "contactid", contactId)
		//r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

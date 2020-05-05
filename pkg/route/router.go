package route

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//ContactController interface
type ContactController interface {
	GetOneContactCtrl(w http.ResponseWriter, r *http.Request)
	AddOrUpdateContactCtrl(w http.ResponseWriter, r *http.Request)
}

type routeFrame struct {
	Name        string
	Method      string
	Pattern     string
	Queries     []string
	HandlerFunc http.HandlerFunc
}

type routeLst []routeFrame

//NewRouter register the routeFrame and handler
func NewRouter(cc ContactController) *mux.Router {
	var routes = routeLst{
		routeFrame{
			"GetOneContactCtrl",
			strings.ToUpper("Get"),
			//"GET", 127.0.0.1:8080/v1/contact/contact_id_or_email
			"/v1/contact/{contact_id_or_email}",
			[]string{},
			cc.GetOneContactCtrl,
		},
		routeFrame{
			"AddOrUpdateContactCtrl",
			strings.ToUpper("Post"),
			//127.0.0.1:8080/v1/contact
			"/v1/contact",
			[]string{},
			cc.AddOrUpdateContactCtrl,
		},
		/*routeFrame{
			"GetAllContactsCtrl",
			strings.ToUpper("Get"),
			//"GET", 127.0.0.1:8080/v1/contacts
			"/v1/contacts",
			[]string{},
			controller.GetAllContactsCtrl,
		},
		routeFrame{
			"GetAllContactsBookmarkCtrl",
			strings.ToUpper("Get"),
			//"GET", 127.0.0.1:8080/v1/contacts/bookmark
			"/v1/contacts/bookmark",
			[]string{},
			controller.GetAllContactsBookmarkCtrl,
		},
		routeFrame{
			"AddBulkContactsCtrl",
			strings.ToUpper("Post"),
			//127.0.0.1:8080/v1/contact
			"/v1/contacts",
			[]string{},
			controller.AddBulkContactsCtrl,
		},
		routeFrame{
			"DeleteContactCtrl",
			strings.ToUpper("Delete"),
			//127.0.0.1:8080/v1/contact
			"/v1/contact/{contact_id_or_email}",
			[]string{},
			controller.DeleteContactCtrl,
		},
		routeFrame{
			"UnsubscribeContactCtrl",
			strings.ToUpper("Post"),
			//127.0.0.1:8080/v1/contact/contact_id_or_email/unsubscribe
			"/v1/contact/{contact_id_or_email}/unsubscribe",
			[]string{},
			controller.UnsubscribeContactCtrl,
		},*/
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Queries(route.Queries...).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		logrus.Debugf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

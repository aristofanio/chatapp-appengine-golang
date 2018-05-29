package adm

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

//--------------------------------------------------------------------
// Routing
//--------------------------------------------------------------------

//create login router
func indexRouter() *mux.Router {
	//create route
	r := mux.NewRouter()
	//mapping handler
	r.HandleFunc("/pub/index.html", pubIndexAuthHandler).Methods("GET")
	r.Handle("/assets/img/splash.png", http.FileServer(http.Dir("/assets"))).Methods("GET")
	//result
	return r
}

//--------------------------------------------------------------------
// Handlers
//--------------------------------------------------------------------

//check if has error
func hasError(err error, w http.ResponseWriter) bool {
	if err != nil {
		fmt.Fprintf(w, "{\"success\": false, \"message\": \"%s\"}", url.QueryEscape(err.Error()))
		return true
	}
	return false
}

//result login page
func pubIndexAuthHandler(w http.ResponseWriter, r *http.Request) {
	//resulting
	tmpl, err := template.ParseFiles("www/pages/pub/index.html")
	if hasError(err, w) {
		return
	}
	//execute template and write in response
	tmpl.Execute(w, nil)
}

//initialize configuration
func init() {
	http.Handle("/pub/", indexRouter()) //public
}

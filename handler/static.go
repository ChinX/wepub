package handler

import (
	"html/template"
	"net/http"

	"github.com/chinx/wepub/module"
)

var StaticDir = "./static"

func EditHandler(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
	t, err := template.ParseFiles(StaticDir + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = urlParam(r, "filename")
	http.FileServer(http.Dir(StaticDir)).ServeHTTP(w, r)
}

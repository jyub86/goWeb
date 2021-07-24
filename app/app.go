package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

type AppHandler struct {
	http.Handler
}

var rd *render.Render = render.New()

func (a *AppHandler) pathHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	imgs, err := Search(path)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	rd.JSON(w, http.StatusOK, map[string][]ImageData{"data": imgs})
}

func MakeHandler() *AppHandler {
	r := mux.NewRouter()
	n := negroni.Classic()
	n.UseHandler(r)
	a := &AppHandler{
		Handler: n,
	}
	r.HandleFunc("/", a.pathHandler).Methods("POST")
	return a
}

package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Auth!\n")
}

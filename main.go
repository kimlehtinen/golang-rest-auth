package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/controllers"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome home!\n")
}

func HelloName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
}

func main() {
	fmt.Println("Hello, golang-rest-auth")

	router := httprouter.New()
	router.GET("/", Home)
	router.GET("/auth", controllers.Auth)
	router.GET("/hello/:name", HelloName)
	log.Fatal(http.ListenAndServe(":8080", router))
}

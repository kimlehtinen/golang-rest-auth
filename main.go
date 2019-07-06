package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/controllers"
	"github.com/kim3z/golang-rest-auth/middleware"
)

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome home!\n")
}

func Foo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome bar!\n")
}

func HelloName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
}

func main() {
	fmt.Println("Hello, golang-rest-auth")

	router := httprouter.New()
	router.GET("/", Home)
	router.GET("/hello/:name", middleware.Authenticate(httprouter.Handle(HelloName)))
	router.GET("/foo", middleware.Authenticate(httprouter.Handle(Foo)))

	router.POST("/api/user/create", controllers.CreateUser)
	router.POST("/api/user/login", controllers.LoginUser)
	router.POST("/api/user/forgot-password/:email", controllers.ForgotPassword)

	log.Fatal(http.ListenAndServe(":8080", router))
}

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

func main() {
	router := httprouter.New()
	router.GET("/", Home)
	router.POST("/api/user/create", controllers.CreateUser)
	router.POST("/api/user/login", controllers.LoginUser)
	router.POST("/api/user/forgot-password/:email", controllers.ForgotPassword)

	// Protected route example
	// router.GET("/foo", middleware.Authenticate(httprouter.Handle(Foo)))

	log.Fatal(http.ListenAndServe(":8080", router))
}

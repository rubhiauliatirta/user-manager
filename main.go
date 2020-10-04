package main

import (
	"fmt"
	"log"
	"net/http"

	_ "./config"
	"./controllers"
	"./helpers"
)

type MyMux struct{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		controllers.Login(w, r)
	case "/register":
		controllers.Register(w, r)
	case "/logout":
		controllers.Logout(w, r)
	default:
		authUser, err := helpers.Auth(r)
		if err != nil {
			http.Redirect(w, r, "/login?auth=false", http.StatusSeeOther)
			return
		}
		fmt.Println(authUser)
		switch r.URL.Path {
		case "/":
			controllers.Home(w, r)
		case "/detail":
			controllers.Detail(w, r)
		case "/edit":
			isSuccess := helpers.IsAuthorized(authUser, r)
			if isSuccess {
				controllers.Edit(w, r)
			} else {
				http.Redirect(w, r, "/?authorize=false", http.StatusSeeOther)
			}
		case "/delete":
			isSuccess := helpers.IsAuthorized(authUser, r)
			if isSuccess {
				controllers.Delete(w, r)
			} else {
				http.Redirect(w, r, "/?authorize=false", http.StatusSeeOther)
			}
		default:
			http.NotFound(w, r)
		}
	}

	return
}

func main() {

	mux := &MyMux{}

	err := http.ListenAndServe(":9090", mux)

	if err != nil {
		log.Fatal("Something error")
	}
}

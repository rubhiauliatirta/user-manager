package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"../config"
	"../models"
)

type MessageData struct {
	IsError bool
	Message string
}

func Login(w http.ResponseWriter, r *http.Request) {

	session, err := config.Store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		tmp, _ := template.ParseFiles("./templates/login.html")

		regValue := r.URL.Query().Get("reg")
		auth := r.URL.Query().Get("auth")

		if regValue == "success" {
			tmp.Execute(w, MessageData{false, "Register success! Please Login"})
		} else if auth == "false" {
			tmp.Execute(w, MessageData{true, "Please login first"})
		} else {
			tmp.Execute(w, nil)
		}

	case "POST":
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		loggedUser, err := models.GetUserByUsername(username, password)

		switch err {
		case sql.ErrNoRows:

			tmp, _ := template.ParseFiles("./templates/login.html")
			tmp.Execute(w, MessageData{true, "Invalid email/password!"})

			return
		case nil:

			authUser := &models.AuthUser{
				IsAuthenticated: true,
				User:            loggedUser,
			}
			session.Values["user"] = authUser

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			fmt.Println(err.Error())
			tmp, _ := template.ParseFiles("./templates/login.html")
			tmp.Execute(w, MessageData{true, "Something Error! Please try again later."})
		}
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmp, _ := template.ParseFiles("./templates/register.html")
		tmp.Execute(w, nil)
	case "POST":
		r.ParseForm()
		newUser, err := models.ValidateUser(r.Form, false)

		if err != nil {
			tmp, _ := template.ParseFiles("./templates/register.html")
			tmp.Execute(w, MessageData{true, err.Error()})

			return
		}

		id, _ := models.RegisterUser(newUser)
		newUser.Id = id

		http.Redirect(w, r, "/login?reg=success", http.StatusSeeOther)

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, "auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = models.AuthUser{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func Home(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	tmp, _ := template.ParseFiles("./templates/home.html")
	edit := r.URL.Query().Get("edit")
	delete := r.URL.Query().Get("delete")
	authorize := r.URL.Query().Get("authorize")
	message := MessageData{false, ""}

	if edit == "success" {
		message.Message = "Edit success!"
	}

	if delete == "success" {
		message.Message = "Delete success!"
	}
	if authorize == "false" {
		message.Message = "You are not authorize to perform the action"
		message.IsError = true
	}

	if err != nil {
		tmp.Execute(w, MessageData{true, err.Error()})
		return
	}
	fmt.Println(users)
	err = tmp.Execute(w, struct {
		Users []models.User
		MessageData
	}{users, message})
	fmt.Println(err)

}

func Detail(w http.ResponseWriter, r *http.Request) {
	detailId := r.URL.Query().Get("id")
	id, err := strconv.Atoi(detailId)
	tmpHome, _ := template.ParseFiles("./templates/home.html")
	tmpDetail, _ := template.ParseFiles("./templates/detail.html")
	if err != nil {
		tmpHome.Execute(w, MessageData{true, err.Error()})
	}
	user, err := models.GetUserById(id)

	if err == sql.ErrNoRows {
		tmpHome.Execute(w, MessageData{true, "User not found!"})
	} else if err == nil {
		tmpDetail.Execute(w, user)
	} else {
		fmt.Println(err)
		tmpHome.Execute(w, MessageData{true, "Something Error!"})
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	editedId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(editedId)
	tmpEdit, _ := template.ParseFiles("./templates/edit.html")
	tmpHome, _ := template.ParseFiles("./templates/home.html")
	switch r.Method {
	case "GET":
		editedUser, err := models.GetUserById(id)
		if err != nil {
			tmpHome.Execute(w, MessageData{true, err.Error()})
			return
		}
		tmpEdit.Execute(w, editedUser)
	case "POST":
		r.ParseForm()

		editedUser, err := models.ValidateUser(r.Form, true)
		editedUser.Id = id

		if err != nil {
			tmpHome.Execute(w, MessageData{true, err.Error()})
			return
		}

		err = models.UpdateUser(editedUser)
		if err != nil {
			tmpHome.Execute(w, MessageData{true, err.Error()})
			return
		}

		http.Redirect(w, r, "/?edit=success", http.StatusSeeOther)

	}
}
func Delete(w http.ResponseWriter, r *http.Request) {
	deleteId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(deleteId)
	tmpHome, _ := template.ParseFiles("./templates/home.html")

	err := models.DeleteUser(id)
	if err != nil {
		tmpHome.Execute(w, MessageData{true, err.Error()})
		return
	}
	http.Redirect(w, r, "/?delete=success", http.StatusSeeOther)
}

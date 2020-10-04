package models

import (
	"encoding/gob"
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"../config"
)

type User struct {
	Id       int
	Username string
	Password string
	Name     string
	Age      int
	Email    string
	Role     string
}

type AuthUser struct {
	IsAuthenticated bool
	User            User
}

func init() {
	gob.Register(AuthUser{})
}

func GetUsers() ([]User, error) {

	users := make([]User, 0)
	rows, err := config.Db.Query("SELECT id, email, name, username, age, role FROM users")
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Name,
			&user.Username,
			&user.Age,
			&user.Role,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil

}

func GetUserByUsername(username string, password string) (User, error) {

	var userResult User
	sqlStatement := `
        SELECT *
        FROM users
        WHERE username=$1 and password=$2
        LIMIT 1;`
	row := config.Db.QueryRow(sqlStatement, username, password)
	err := row.Scan(
		&userResult.Id,
		&userResult.Username,
		&userResult.Password,
		&userResult.Age,
		&userResult.Email,
		&userResult.Name,
		&userResult.Role,
	)

	return userResult, err

}

func RegisterUser(newUser User) (int, bool) {

	sqlStatement := `
		INSERT INTO users (username, password, email, name, age, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int
	row := config.Db.QueryRow(sqlStatement,
		newUser.Username,
		newUser.Password,
		newUser.Email,
		newUser.Name,
		newUser.Age,
		"user")

	err := row.Scan(&id)

	switch err {
	case nil:
		return id, true
	default:
		panic(err)
	}
}

func GetUserById(id int) (User, error) {
	var userResult User
	sqlStatement := `
        SELECT *
        FROM users
        WHERE id=$1
        LIMIT 1;`
	row := config.Db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&userResult.Id,
		&userResult.Username,
		&userResult.Password,
		&userResult.Age,
		&userResult.Name,
		&userResult.Email,
		&userResult.Role,
	)
	return userResult, err

}

func UpdateUser(updatedUser User) error {
	sqlStatement := `
		UPDATE users
		SET username=$2, name=$3, email=$4, age=$5
		WHERE id=$1
	`
	_, err := config.Db.Exec(
		sqlStatement,
		updatedUser.Id,
		updatedUser.Username,
		updatedUser.Name,
		updatedUser.Email,
		updatedUser.Age,
	)
	if err != nil {

		return err
	}

	return nil

}

func DeleteUser(id int) error {
	sqlStatement := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := config.Db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil

}

func ValidateUser(val url.Values, isUpdate bool) (User, error) {
	username := val.Get("username")
	password := val.Get("password")
	name := val.Get("name")
	age := val.Get("age")
	email := val.Get("email")

	var newUser User
	var validationError error = nil
	errorStrings := make([]string, 0, 4)

	if len(username) < 5 {
		errorStrings = append(errorStrings, "Mininum username length is 5")
	}

	if len(password) < 8 && !isUpdate {
		errorStrings = append(errorStrings, "Minimal karakter password adalah 8")
	}
	m, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)

	if !m {
		errorStrings = append(errorStrings, "Format email salah")
	}

	valInt, err := strconv.Atoi(age)
	if err != nil || valInt < 0 {
		errorStrings = append(errorStrings, "Age bukan bilangan bulat positif")
	}

	if len(errorStrings) != 0 {
		validationError = errors.New(strings.Join(errorStrings, ","))
	}

	if validationError == nil {
		newUser = User{
			Name:     name,
			Username: username,
			Password: password,
			Age:      valInt,
			Email:    email,
		}
	}
	return newUser, validationError

}

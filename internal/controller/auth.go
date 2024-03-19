package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"forum/internal/model"
)

// Login, done
func Login(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("login")))
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}

// SignUp
func (ctl *BaseController) SignUp(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("login")))

	// receive user data from formValue
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// validation data, implemented in frontent side

	fmt.Println(username, email, password)
	// if user with this email already exists, return error
	usr, err := ctl.Repo.URepo.GetUserByEmail(email)
	if usr != nil {
		slog.Error(err.Error())
		errMsg := "User with this email, allready exists. Please, try again with another email!"
		tmp.Execute(w, errMsg)
		return
	}

	// hashing password
	hash, err := GetHashFromPassword(password)
	if err != nil {
		slog.Error(err.Error())
		errMsg := "Invalid password, enter correct password, and try again!"
		tmp.Execute(w, errMsg)
		return
	}

	// create a new user instance with this data
	user := &model.User{
		Username: username,
		Email:    email,
		Password: hash,
	}

	// insert this user into database
	userID, err := ctl.Repo.URepo.CreateUser(user)
	if err != nil {
		slog.Error(err.Error())
		errMsg := "Invalid Data. Please, try again!"
		tmp.Execute(w, errMsg)
		return
	}

	// create a new cookie, and set this cookie to user
	uuid := uuid.DefaultGenerator
	suuid, _ := uuid.NewV4()
	sValue := suuid.String()

	expires := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sValue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}

	// CreateSession -> userID, sessionID, expires
	if err := ctl.Repo.URepo.CreateSession(sValue, userID, expires); err != nil {
		slog.Error(err.Error())
		errMsg := "Invalid Data. Please, try again!"
		tmp.Execute(w, errMsg)
		return
	}

	http.SetCookie(w, cookie)

	// redirect to main page
	fmt.Println("SESSIONNNNN", cookie, "USERID:", userID)
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
}

func Logout(w http.ResponseWriter, r *http.Request) {
}

// GetHashFromPassword
func GetHashFromPassword(password string) (passwordHash string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"forum/internal/model"
)

// GetWD
func GetWD() (wd string) {
	wd, _ = os.Getwd()
	return wd
}

// Login, done
func Login(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("login")))
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}

// Sign-up page
func (ctl *BaseController) SignUpPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_up.html"))
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}

// Login, done
func (ctl *BaseController) SignInPage(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_in.html"))
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}

// SignUp
func (ctl *BaseController) SignUp(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_up.html"))

	// receive user data from formValue
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// validation data, implemented in frontent side

	fmt.Println(username, email, password)
	// if user with this email already exists, return error
	us := ctl.Repo.URepo.GetUserByEmail(email)
	if us != nil {
		fmt.Println("USER BEFORE:", us)
		errMsg := "User with this email, allready exists. Please, try again with another email!"
		tmp.Execute(w, errMsg)
		return
	}

	// hashing password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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
		Password: string(hash),
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
	ids := strconv.Itoa(userID)
	svalue := strings.Join([]string{ids, sValue}, ",")

	expires := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    svalue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}

	// CreateSession -> userID, sessionID, expires
	if err := ctl.Repo.URepo.CreateSession(userID, expires); err != nil {
		slog.Error(err.Error())
		errMsg := "Invalid Data. Please, try again!"
		tmp.Execute(w, errMsg)
		return
	}

	http.SetCookie(w, cookie)

	// redirect to main page
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

func (ctl *BaseController) SignIn(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_in.html"))
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := ctl.Repo.URepo.GetUserByEmail(email)
	if user == nil {
		errMsg := "Invalid Email, try again!"
		tmpl.Execute(w, errMsg)
		return
	}

	fmt.Println("User hash", user.Password)
	fmt.Println("User ID", user.ID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("\n\n\nERRROR", err.Error())
		errMsg := "Invalid Password. Please, try again!"
		tmpl.Execute(w, errMsg)
		return
	}

	// create a new cookie, and set this cookie to user
	uuid := uuid.DefaultGenerator
	suuid, _ := uuid.NewV4()
	sValue := suuid.String()

	ids := strconv.Itoa(user.ID)
	svalue := strings.Join([]string{ids, sValue}, ",")

	expires := time.Now().Add(1 * time.Hour)
	coockie := &http.Cookie{
		Name:     "sessionID",
		Value:    svalue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}

	fmt.Println("COOCKIE:", coockie)
	http.SetCookie(w, coockie)
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
	return
}

func (ctl *BaseController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
	}

	// get userID from cookie
	sValue := cookie.Value
	svalue := strings.Split(sValue, ",")

	userID, err := strconv.Atoi(svalue[0])
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
	}

	if err := ctl.Repo.URepo.RemoveSession(userID); err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
	}

	slog.Info("User logout")
	newCookie := &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, newCookie)

	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

func (ctl *BaseController) GetAuthUser(r *http.Request) (user *model.User) {
	// get userID from cookie
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return nil
	}
	sValue := cookie.Value
	svalue := strings.Split(sValue, ",")

	userID, err := strconv.Atoi(svalue[0])
	if err != nil {
		return nil
	}

	user, err = ctl.Repo.URepo.GetUserByID(userID)
	if err != nil {
		return nil
	}

	return user
}

// GetCockie
func (ctl *BaseController) GetCoockie(user *model.User) *http.Cookie {
	// create a new cookie, and set this cookie to user
	uuid := uuid.DefaultGenerator
	suuid, _ := uuid.NewV4()
	sValue := suuid.String()

	expires := time.Now().Add(1 * time.Hour)
	coockie := &http.Cookie{
		Name:     "sessionID",
		Value:    sValue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}
	// user = ctl.Repo.URepo.GetUserByEmail()

	// ctl.Repo.URepo.RemoveSession(user.ID)

	ctl.Repo.URepo.CreateSession(user.ID, expires)

	return coockie
}

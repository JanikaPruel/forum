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

// SignUpPage - This function is a method of the BaseController struct, and it
// serves the sign-up page to the user. It takes in a ResponseWriter and a Request,
// and it uses the ParseFiles function from the template package to parse the
// sign-up page template file. It then sets the status code of the ResponseWriter
// to 200 OK, and executes the template, passing in nil as the data. If there
// is an error during template execution, it logs the error.
func (ctl *BaseController) SignUpPage(w http.ResponseWriter, r *http.Request) {
	// Parse the sign-up page template file
	tmp := template.Must(template.ParseFiles(
		GetWD() + "/web/templates/sign_up.html"))

	// Set the status code of the ResponseWriter to 200 OK
	w.WriteHeader(http.StatusOK)

	// Execute the template, passing in nil as the data
	if err := tmp.Execute(w, nil); err != nil {
		// If there is an error, log it
		slog.Error(err.Error())
	}
}

// SignInPage renders the sign-in page template to the ResponseWriter.
//
// This function is a method of the BaseController struct. It takes in a
// http.ResponseWriter and a http.Request as parameters. It uses the ParseFiles
// function from the template package to parse the sign-in page template file.
// It then sets the status code of the ResponseWriter to 200 OK, and executes
// the template, passing in nil as the data. If there is an error during
// template execution, it logs the error.
func (ctl *BaseController) SignInPage(w http.ResponseWriter, r *http.Request) {
	// Parse the sign-in page template file
	tmp := template.Must(template.ParseFiles(
		GetWD() + "/web/templates/sign_in.html"))

	// Set the status code of the ResponseWriter to 200 OK
	w.WriteHeader(http.StatusOK)

	// Execute the template, passing in nil as the data
	if err := tmp.Execute(w, nil); err != nil {
		// If there is an error, log it
		slog.Error(err.Error())
	}

}

// SignUp handles the user sign-up process.
//
// This function receives user data from the formValue, validates it,
// hashes the password, creates a new user, inserts the user into the
// database, creates a new session and sets a cookie for the user.
// If any error occurs during the process, it logs the error and
// renders an error message to the user. Finally, it redirects the
// user to the main page.
func (ctl *BaseController) SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse the sign-up page template file
	tmp := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_up.html"))

	// Receive user data from formValue
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate data (implemented in front-end side)
	fmt.Println(username, email, password)

	// Check if user with this email already exists, return error
	us := ctl.Repo.URepo.GetUserByEmail(email)
	if us != nil {
		// User with this email already exists, return error
		fmt.Println("USER BEFORE:", us)
		errMsg := "User with this email, allready exists. Please, try again with another email!"
		tmp.Execute(w, errMsg)
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		// Invalid password, return error
		slog.Error(err.Error())
		errMsg := "Invalid password, enter correct password, and try again!"
		tmp.Execute(w, errMsg)
		return
	}

	// Create a new user instance with this data
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hash),
	}

	// Insert this user into database
	userID, err := ctl.Repo.URepo.CreateUser(user)
	if err != nil {
		// Invalid Data, return error
		slog.Error(err.Error())
		errMsg := "Invalid Data. Please, try again!"
		tmp.Execute(w, errMsg)
		return
	}

	// Create a new cookie and set it to the user
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

	// Create session in the database
	if err := ctl.Repo.URepo.CreateSession(userID, expires); err != nil {
		slog.Error(err.Error())
		errMsg := "Invalid Data. Please, try again!"
		tmp.Execute(w, errMsg)
		return
	}

	http.SetCookie(w, cookie)

	// Redirect to main page
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

// SignIn handles the user sign-in process. It parses the HTML template for the sign-in page,
// retrieves the user's email and password from the request form, and verifies the user's
// credentials. If the credentials are valid, it creates a new cookie and sets it to the user.
// The cookie contains the user's ID and a generated UUID. The function then redirects the user
// to the home page.
func (ctl *BaseController) SignIn(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template for the sign-in page
	tmpl := template.Must(template.ParseFiles(GetWD() + "/web/templates/sign_in.html"))

	// Retrieve the user's email and password from the request form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Retrieve the user from the database based on the email
	user := ctl.Repo.URepo.GetUserByEmail(email)

	// If the user does not exist, display an error message and return
	if user == nil {
		errMsg := "Invalid Email, try again!"
		tmpl.Execute(w, errMsg)
		return
	}

	// Print the user's hashed password and ID for debugging purposes
	fmt.Println("User hash", user.Password)
	fmt.Println("User ID", user.ID)

	// Verify the user's password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("\n\n\nERRROR", err.Error())
		errMsg := "Invalid Password. Please, try again!"
		tmpl.Execute(w, errMsg)
		return
	}

	// Generate a new UUID and join it with the user's ID to create a cookie value
	uuid := uuid.DefaultGenerator
	suuid, _ := uuid.NewV4()
	sValue := suuid.String()
	ids := strconv.Itoa(user.ID)
	svalue := strings.Join([]string{ids, sValue}, ",")

	// Set the cookie expiration to 1 hour
	expires := time.Now().Add(1 * time.Hour)

	// Create a new cookie and set it to the user
	coockie := &http.Cookie{
		Name:     "sessionID",
		Value:    svalue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}

	// Print the cookie for debugging purposes
	fmt.Println("COOCKIE:", coockie)

	// Set the cookie in the response and redirect the user to the home page
	http.SetCookie(w, coockie)
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

// Logout handles the user logout process. It removes the sessionID cookie and
// redirects the user to the home page.
func (ctl *BaseController) Logout(w http.ResponseWriter, r *http.Request) {
	// Retrieve the sessionID cookie from the request
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// If the cookie is not found, log the error and redirect to the home page with 401 status
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
		return
	}

	// Get the userID from the cookie
	sValue := cookie.Value
	svalue := strings.Split(sValue, ",")

	userID, err := strconv.Atoi(svalue[0])
	if err != nil {
		// If the userID cannot be retrieved from the cookie, log the error and redirect to the home page with 401 status
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
		return
	}

	// Remove the session for the user
	if err := ctl.Repo.URepo.RemoveSession(userID); err != nil {
		// If the session removal fails, log the error and redirect to the home page with 401 status
		slog.Error(err.Error())
		http.Redirect(w, r, "GET /", http.StatusUnauthorized)
		return
	}

	// Log the user logout
	slog.Info("User logout")

	// Create a new cookie with empty value and expiration in the past
	newCookie := &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Unix(0, 0),
	}

	// Set the new cookie in the response and redirect the user to the home page
	http.SetCookie(w, newCookie)
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

// GetAuthUser retrieves the authenticated user from the request.
// It does this by extracting the user ID from the "sessionID" cookie,
// retrieving the corresponding user from the user repository, and returning it.
// If any error occurs during this process, it returns nil.
func (ctl *BaseController) GetAuthUser(r *http.Request) (user *model.User) {
	// Extract the user ID from the "sessionID" cookie.
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// If the cookie is not found, return nil.
		return nil
	}

	// Split the cookie value by comma.
	sValue := cookie.Value
	svalue := strings.Split(sValue, ",")

	// Convert the first part of the split value to an integer.
	userID, err := strconv.Atoi(svalue[0])
	if err != nil {
		// If the conversion fails, return nil.
		return nil
	}

	// Retrieve the user from the user repository using the user ID.
	user, err = ctl.Repo.URepo.GetUserByID(userID)
	if err != nil {
		// If any error occurs during retrieval, return nil.
		return nil
	}

	// Return the retrieved user.
	return user
}

// GetCoockie generates and returns a new HTTP cookie for the given user.
//
// The cookie is named "sessionID" and its value is a new UUID v4 generated by the
// uuid package. The expiration time of the cookie is set to 1 hour from the
// current time. The cookie is set to be HTTP only and secure.
//
// The function also creates a new session in the user repository for the given
// user, using the session duration of 1 hour.
//
// Parameters:
// - user: The user for whom the cookie is generated.
//
// Returns:
// The generated HTTP cookie.
func (ctl *BaseController) GetCoockie(user *model.User) *http.Cookie {
	// Generate a new UUID v4 for the cookie value.
	uuid := uuid.DefaultGenerator
	suuid, _ := uuid.NewV4()
	sValue := suuid.String()

	// Set the expiration time of the cookie to 1 hour from the current time.
	expires := time.Now().Add(1 * time.Hour)

	// Create a new HTTP cookie with the generated value and expiration time.
	coockie := &http.Cookie{
		Name:     "sessionID",
		Value:    sValue,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}

	// Create a new session in the user repository for the given user.
	ctl.Repo.URepo.CreateSession(user.ID, expires)

	// Return the generated cookie.
	return coockie
}

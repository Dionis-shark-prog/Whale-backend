package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var hashKeyUser = []byte("very-secret")
var blockKeyUser = []byte("a-lot-secret-key")
var s2 = securecookie.New(hashKeyUser, blockKeyUser)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	redirected := r.URL.Query().Get("redirected")
	var message string
	if redirected == "cart" {
		message = "Log in to use Cart and buy stuff"
	} else if redirected == "logout" || redirected == "delete" {
		message = "You are not loged in"
	} else if redirected == "profile" {
		message = "Something wrong when opening profile"
	}

	login := "login.html"

	loginContent, err := getresources.GetContentByURL(login, "login")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("login").Parse(loginContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstname := r.FormValue("firstname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println("firstname:", firstname, "email:", email, "password:", password)

	hash, err := models.TakeHashByNameEmail(firstname, email)
	if err != nil {
		login := "login.html"
		message := err.Error()

		loginContent, err := getresources.GetContentByURL(login, "login")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("login").Parse(loginContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		cookie, err := login(firstname, email)
		if err != nil {
			login := "login.html"
			message := err.Error()

			loginContent, err := getresources.GetContentByURL(login, "login")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.New("login").Parse(loginContent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, fmt.Sprint("error while loging in:", message))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/shop?in=true", http.StatusSeeOther)
		return
	}
	fmt.Println("Incorrect password!")

	signup := "../web/templates/login.html"
	tmpl := template.Must(template.ParseFiles(signup))
	tmpl.Execute(w, "Incorrect password!")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieUser, err := r.Cookie("user_seccion")
	if err != nil {

		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	userData, err := DecodeUserCookie(cookieUser.Value)
	if err != nil {

		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	_, err = models.TakeUserByToken(userData["user_token"])
	if err != nil {

		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	cookieUser.MaxAge = -1
	http.SetCookie(w, cookieUser)

	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	cookieUser, err := r.Cookie("user_seccion")
	if err != nil {

		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	userData, err := DecodeUserCookie(cookieUser.Value)
	if err != nil {

		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	id, err := models.TakeUserByToken(userData["user_token"])
	if err != nil {

		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=logout", http.StatusSeeOther)
		return
	}

	err = models.DeleteClientFromBase(id)
	if err != nil {
		fmt.Println(w, err.Error())
		return
	}

	cookieUser.MaxAge = -1
	http.SetCookie(w, cookieUser)

	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func login(firstname string, email string) (*http.Cookie, error) {
	userToken, err := models.TakeTokenByNameEmail(firstname, email)
	if err != nil {
		return nil, err
	}

	value := map[string]string{
		"user_token": userToken,
	}

	encoded, err := s2.Encode("user_seccion", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:     "user_seccion",
			Value:    encoded,
			MaxAge:   3600,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		return cookie, nil
	}
	return nil, err
}

func DecodeUserCookie(code string) (map[string]string, error) {
	value := make(map[string]string)
	err := s1.Decode("user_seccion", code, &value)
	if err == nil {
		return value, nil
	}
	return nil, err
}

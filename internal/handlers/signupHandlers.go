package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	signup := "signup.html"

	signupContent, err := getresources.GetContentByURL(signup, "signup")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("signup").Parse(signupContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignupComplete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstname := r.FormValue("firstname")

	// 1. Checking if name is correct
	firstnameAlphaNumeric := firstname != ""
	for _, char := range firstname {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			firstnameAlphaNumeric = false
			break
		}
	}

	firstnameLength := false
	if 2 <= len(firstname) && len(firstname) <= 30 {
		firstnameLength = true
	}

	checkFirstName := firstnameAlphaNumeric && firstnameLength

	secondname := r.FormValue("secondname")

	// 2. Checking if secondname is correct
	secondnameAlphaNumeric := true
	for _, char := range secondname {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			secondnameAlphaNumeric = false
		}
	}

	secondnameLength := false
	if len(secondname) <= 30 {
		secondnameLength = true
	}

	checkSecondName := secondnameAlphaNumeric && secondnameLength

	email := r.FormValue("email")

	// 3. Checking if email is correct
	emailLength := false
	if 0 < len(email) && len(email) <= 60 {
		emailLength = true
	}

	checkEmail := emailLength

	password := r.FormValue("password")

	// 4. Checking if password is correct
	pswdIsEmpty := password != ""
	pswdHasSpaces := regexp.MustCompile(`\s`).MatchString(password)

	checkPassword := pswdIsEmpty && !pswdHasSpaces

	number := r.FormValue("number")

	// 5. Checking if number is correct
	numberHasNumerics := true
	for _, char := range number {
		if !unicode.IsNumber(char) {
			numberHasNumerics = false
			break
		}
	}

	checkNumber := numberHasNumerics

	if !checkFirstName || !checkSecondName || !checkEmail || !checkPassword || !checkNumber {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "Check your data: missing mandatory fields or incorect input")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	status, err := models.IsUserByName(firstname, email)

	if status == 0 {
		signup := "signup.html"
		message := err.Error()

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if status == 2 {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "User already exists")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var hash []byte

	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "A problem with registering an account")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var secondnamePtr, numberPtr *string
	if secondname != "" {
		secondnamePtr = &secondname
	}
	if number != "" {
		numberPtr = &number
	}

	err = models.InsertUserToClients(firstname, secondnamePtr, email, hash, numberPtr)
	if err != nil {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "Error while inserting data!")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	idForToken, err := models.TakeIDByNameEmail(firstname, email)
	if err != nil {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "Error while trying to create token!")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = models.InsertAndCreateToken(idForToken)
	if err != nil {
		signup := "signup.html"

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, "Error while inserting into tokens!")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	cookie, err := login(firstname, email)
	if err != nil {
		signup := "signup.html"
		message := err

		signupContent, err := getresources.GetContentByURL(signup, "signup")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("signup").Parse(signupContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, fmt.Sprint("Error while loging in:", message))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/shop?in=true", http.StatusSeeOther)
}

package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, UserID := CheckIfUserCookieIsCorrect(w, r)

	userData, err := models.TakeUserDataByID(UserID)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login?redirected=profile", http.StatusSeeOther)
		return
	} else if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	profile := "profile.html"

	adminContent, err := getresources.GetContentByURL(profile, "profile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("profile").Parse(adminContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

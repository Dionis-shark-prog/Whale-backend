package middleware

import (
	"WhaleWebSite/internal/handlers"
	"WhaleWebSite/internal/models"
	"context"
	"net/http"
)

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cUser, err := r.Cookie("user_seccion")
		isAuthenticated := false

		if err == nil && cUser.Value != "" {
			isAuthenticated = true
		}

		ctx := context.WithValue(r.Context(), "IsAuthenticated", isAuthenticated)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CheckIfUserCookieIsCorrect(w http.ResponseWriter, r *http.Request) (*http.Cookie, int64) {
	cookieUser, err := r.Cookie("user_seccion")
	if err != nil {
		http.Redirect(w, r, "/login?redirected=cart", http.StatusSeeOther)
		return nil, 0
	}

	userData, err := handlers.DecodeUserCookie(cookieUser.Value)
	if err != nil {
		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=cart", http.StatusSeeOther)
		return nil, 0
	}

	userID, err := models.TakeUserByToken(userData["user_token"])
	if err != nil {
		cookieUser.MaxAge = -1
		http.SetCookie(w, cookieUser)
		http.Redirect(w, r, "/login?redirected=cart", http.StatusSeeOther)
		return nil, 0
	}
	return cookieUser, userID
}

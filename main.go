package main

import (
	"WhaleWebSite/internal/handlers"
	"WhaleWebSite/internal/models"
	"WhaleWebSite/middleware"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	models.DBConnect()

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:8080/static"+r.URL.Path[len("/static/"):], http.StatusFound)
	})

	mux.HandleFunc("/images/", handlers.ImageHandler)

	mux.HandleFunc("/shop", handlers.ShopHandler)
	mux.HandleFunc("/shop/goods", handlers.GoodsHandler)

	mux.HandleFunc("/about", handlers.AboutHandler)

	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logincomplete", handlers.LoginAuthHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	mux.HandleFunc("/deleteaccount", handlers.DeleteAccountHandler)

	mux.HandleFunc("/signup", handlers.SignupHandler)
	mux.HandleFunc("/signupcomplete", handlers.SignupComplete)

	mux.HandleFunc("/cart", handlers.CartHandler)
	mux.HandleFunc("/cartbuy", handlers.CartBuyHandler)
	mux.HandleFunc("/cartdelete", handlers.CartDeleteHandler)
	mux.HandleFunc("/cartclear", handlers.CartClearHandler)

	mux.HandleFunc("/authadmin", handlers.AuthenticateAdminHandler)
	mux.HandleFunc("/authadmincomplete", handlers.AuthenticateCompleteHandler)
	mux.HandleFunc("/admin", handlers.AdminHandler)
	mux.HandleFunc("/admindelete", handlers.AdminDeleteGoodsHandler)
	mux.HandleFunc("/adminedit", handlers.AdminEditGoodsHandler)
	mux.HandleFunc("/editcomplete", handlers.AdminEditGoodsCompleteHandler)
	mux.HandleFunc("/addgoods", handlers.AddGoodsHandler)
	mux.HandleFunc("/addcomplete", handlers.AddGoodsCompleteHandler)
	mux.HandleFunc("/exitadmin", handlers.ExitAdminHandler)

	mux.HandleFunc("/profile", handlers.ProfileHandler)

	err := http.ListenAndServe("localhost:8000", middleware.AuthMiddleware(mux))
	log.Fatal(err)
}

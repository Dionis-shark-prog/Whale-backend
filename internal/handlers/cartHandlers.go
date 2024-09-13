package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func CartHandler(w http.ResponseWriter, r *http.Request) {
	_, userID := CheckIfUserCookieIsCorrect(w, r)

	cartContent, err := models.GetAllFromCartWithID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productsInCart := ProductsInCart{
		CartContents: cartContent,
	}

	cart := "cart.html"

	cartCode, err := getresources.GetContentByURL(cart, "cart")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("cart").Parse(cartCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, productsInCart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CartBuyHandler(w http.ResponseWriter, r *http.Request) {
	_, userID := CheckIfUserCookieIsCorrect(w, r)

	strID := r.URL.Query().Get("id")

	productID, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	currentProduct, err := models.GetProduct(int64(productID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if currentProduct.Count < 1 {
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}

	isProductInCart := models.IsProductByIDUserID(userID, currentProduct.ID)
	if isProductInCart {
		http.Redirect(w, r, fmt.Sprintf("/shop/goods?number=%s", strID)+"&r=isincart", http.StatusSeeOther)
		return
	}

	err = models.InsertProductToCart(userID, currentProduct.ID, currentProduct.Price)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/shop/goods?number=%s", strID), http.StatusSeeOther)
}

func CartDeleteHandler(w http.ResponseWriter, r *http.Request) {
	_, userID := CheckIfUserCookieIsCorrect(w, r)

	strID := r.URL.Query().Get("g")

	productID, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteProductFromCart(userID, int64(productID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func CartClearHandler(w http.ResponseWriter, r *http.Request) {
	_, userID := CheckIfUserCookieIsCorrect(w, r)

	err := models.ClearUserCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func CheckIfUserCookieIsCorrect(w http.ResponseWriter, r *http.Request) (*http.Cookie, int64) {
	cookieUser, err := r.Cookie("user_seccion")
	if err != nil {
		http.Redirect(w, r, "/login?redirected=cart", http.StatusSeeOther)
		return nil, 0
	}

	userData, err := DecodeUserCookie(cookieUser.Value)
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

func CookieDeleteHandler(w http.ResponseWriter, cookieUser *http.Cookie) {
	cookieUser.MaxAge = -1
	http.SetCookie(w, cookieUser)
}

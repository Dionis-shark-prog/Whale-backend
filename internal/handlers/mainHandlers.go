package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func ShopHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := r.Context().Value("IsAuthenticated").(bool)

	productContent := ProductContent{
		Info:            models.GetAllProducts(),
		IsAuthenticated: isAuthenticated,
	}

	base := "base.html"
	shop := "index.html"

	baseContent, err := getresources.GetContentByURL(base, "base")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shopContent, err := getresources.GetContentByURL(shop, "index")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("base").Parse(baseContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.Parse(shopContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, productContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GoodsHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := r.Context().Value("IsAuthenticated").(bool)
	strID := r.URL.Query().Get("number")

	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusInternalServerError)
		return
	}

	currentProduct, err := models.GetProduct(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productInfo := ProductInfo{
		ProdInfo:        currentProduct,
		IsAuthenticated: isAuthenticated,
	}

	base := "base.html"
	goods := "goods.html"

	baseContent, err := getresources.GetContentByURL(base, "base")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	goodsContent, err := getresources.GetContentByURL(goods, "goods")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("base").Parse(baseContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.Parse(goodsContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, productInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := r.Context().Value("IsAuthenticated").(bool)

	aboutInfo := AboutInfo{
		IsAuthenticated: isAuthenticated,
	}

	base := "base.html"
	about := "about.html"

	baseContent, err := getresources.GetContentByURL(base, "base")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aboutContent, err := getresources.GetContentByURL(about, "about")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("base").Parse(baseContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err = tmpl.Parse(aboutContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, aboutInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

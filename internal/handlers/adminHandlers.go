package handlers

import (
	"WhaleWebSite/internal/models"
	getresources "WhaleWebSite/pkg/getResources"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var hashKeyAdmin = []byte("very-secret")
var blockKeyAdmin = []byte("a-lot-secret-key")
var s1 = securecookie.New(hashKeyAdmin, blockKeyAdmin)

func AuthenticateAdminHandler(w http.ResponseWriter, r *http.Request) {
	admin := "authadmin.html"

	adminContent, err := getresources.GetContentByURL(admin, "authadmin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("authadmin").Parse(adminContent)
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

func AuthenticateCompleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("adminname")
	password := r.FormValue("password")

	hash, err := models.TakeHashByAdminName(name)
	if err != nil {
		admin := "authadmin.html"
		message := err.Error()

		adminContent, err := getresources.GetContentByURL(admin, "authadmin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("authadmin").Parse(adminContent)
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

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		cookie, err := adminCookie(name)
		if err != nil {
			admin := "authadmin.html"
			message := err.Error()

			adminContent, err := getresources.GetContentByURL(admin, "authadmin")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.New("authadmin").Parse(adminContent)
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
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	fmt.Println("Incorrect password!")

	admin := "authadmin.html"

	adminContent, err := getresources.GetContentByURL(admin, "authadmin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("authadmin").Parse(adminContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "Incorrect password!")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	products := ProductAdmin{
		Products: models.GetAllProducts(),
	}

	admin := "admin.html"

	adminContent, err := getresources.GetContentByURL(admin, "admin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("admin").Parse(adminContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AdminDeleteGoodsHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	strID := r.URL.Query().Get("p")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteProduct(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminEditGoodsHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	strID := r.URL.Query().Get("p")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	productToEdit, err := models.GetProduct(int64(id))
	if err != nil {
		http.Error(w, fmt.Sprint("Error while taking data:", err.Error()), http.StatusBadRequest)
		return
	}

	edit := "editgoods.html"

	adminContent, err := getresources.GetContentByURL(edit, "editgoods")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("editgoods").Parse(adminContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, productToEdit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AdminEditGoodsCompleteHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	strID := r.URL.Query().Get("p")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	priceString := r.FormValue("price")
	description := r.FormValue("description")
	countString := r.FormValue("count")

	price, err := strconv.ParseFloat(priceString, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, err := strconv.ParseInt(countString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = models.EditProduct(int64(id), title, description, price, int(count))
	if err != nil {
		http.Error(w, fmt.Sprint("Error while editing data:", err.Error()), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AddGoodsHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	addgoods := "addgoods.html"

	adminContent, err := getresources.GetContentByURL(addgoods, "addgoods")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("addgoods").Parse(adminContent)
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

func AddGoodsCompleteHandler(w http.ResponseWriter, r *http.Request) {
	_ = checkIfAdminCookieIsCorrect(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	price := r.FormValue("price")
	description := r.FormValue("description")
	count := r.FormValue("count")

	files := r.MultipartForm.File["file"]
	fileNames := make([]string, 0, len(files))

	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		imageName := filepath.Join("images", header.Filename) // Incorrect!

		outFile, err := os.Create(imageName)
		if err != nil {
			http.Error(w, "Failed to save image 2!", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Failed to save image 3!", http.StatusInternalServerError)
			return
		}

		fileNames = append(fileNames, header.Filename)
	}

	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		addgoods := "addgoods.html"
		message := err.Error()

		adminContent, err := getresources.GetContentByURL(addgoods, "addgoods")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("addgoods").Parse(adminContent)
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

	countInt, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		addgoods := "addgoods.html"
		message := err.Error()

		adminContent, err := getresources.GetContentByURL(addgoods, "addgoods")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("addgoods").Parse(adminContent)
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

	err = models.InsertGoodsIntoTable(title, description, priceFloat, int(countInt), fileNames)
	if err != nil {
		addgoods := "addgoods.html"
		message := err.Error()

		adminContent, err := getresources.GetContentByURL(addgoods, "addgoods")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("addgoods").Parse(adminContent)
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

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ExitAdminHandler(w http.ResponseWriter, r *http.Request) {
	cookieAdmin := checkIfAdminCookieIsCorrect(w, r)

	cookieAdmin.MaxAge = -1
	http.SetCookie(w, cookieAdmin)

	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

func adminCookie(firstname string) (*http.Cookie, error) {
	id, err := models.TakeIDByAdminName(firstname)
	if err != nil {
		return nil, err
	}

	value := map[string]int64{
		"admin_token": id,
	}

	if encoded, err := s1.Encode("admin_seccion", value); err == nil {
		cookie := &http.Cookie{
			Name:     "admin_seccion",
			Value:    encoded,
			MaxAge:   1800,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		return cookie, nil
	}
	return nil, err
}

func decodeAdminCookie(code string) (map[string]int64, error) {
	value := make(map[string]int64)
	err := s1.Decode("admin_seccion", code, &value)
	if err == nil {
		return value, nil
	}
	return nil, err
}

func checkIfAdminCookieIsCorrect(w http.ResponseWriter, r *http.Request) *http.Cookie {
	cookieAdmin, err := r.Cookie("admin_seccion")
	if err != nil {
		http.Redirect(w, r, "/authadmin", http.StatusSeeOther)
		return nil
	}

	adminData, err := decodeAdminCookie(cookieAdmin.Value)
	if err != nil {
		cookieAdmin.MaxAge = -1
		http.SetCookie(w, cookieAdmin)
		http.Redirect(w, r, "/authadmin", http.StatusSeeOther)
		return nil
	}

	err = models.IsAdminByID(adminData["admin_token"])
	if err != nil {
		cookieAdmin.MaxAge = -1
		http.SetCookie(w, cookieAdmin)
		http.Redirect(w, r, "/authadmin", http.StatusSeeOther)
		return nil
	}
	return cookieAdmin
}

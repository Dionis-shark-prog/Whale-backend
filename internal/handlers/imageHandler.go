package handlers

import (
	"net/http"
	"path/filepath"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	imgUrl := r.URL.Path[len("/images/"):]
	imgPath := filepath.Join("images", imgUrl)

	// buf, err := os.ReadFile(fmt.Sprint("../images", imgName))
	// if err != nil {
	// 	http.Error(w, "Cannot read image content", http.StatusInternalServerError)
	// 	return
	// }

	// http.HandleFunc(w, r, )

	http.ServeFile(w, r, imgPath)
}

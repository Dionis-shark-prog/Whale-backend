package getresources

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetContentByURL(filename string, name string) (string, error) {
	resp, err := http.Get("http://localhost:8080/templates/" + filename)
	if err != nil {
		return "", errors.New(fmt.Sprint("error while getting resource:", err.Error()))
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprint("error while reading content:", err.Error()))
	}

	return string(content), nil
}

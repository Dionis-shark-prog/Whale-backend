package models

import (
	"encoding/json"
	"fmt"
	"log"
)

type Goods struct {
	ID          int64
	Title       string
	Price       float64
	Description string
	Count       int
	ImagesUrls  []string
}

func GetAllProducts() []Goods {
	var goods []Goods

	rows, err := db.Query("SELECT * FROM goods")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var product Goods
		var imageUrl string

		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Description, &product.Count, &imageUrl); err != nil {
			fmt.Println(err.Error())
			return nil
		}

		if err := json.Unmarshal([]byte(imageUrl), &product.ImagesUrls); err != nil {
			log.Fatal(err)
		}

		goods = append(goods, product)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return goods
}

// Function GetProduct(id int64) (Goods, error) searches for a product
// in database by `id` parameter
// If `id` is invalid, `err` will not be nil
func GetProduct(id int64) (Goods, error) {
	var goods Goods
	var imageUrl string

	stmt := "SELECT * FROM goods WHERE ID=?"
	row := db.QueryRow(stmt, id)

	err := row.Scan(&goods.ID, &goods.Title, &goods.Price, &goods.Description, &goods.Count, &imageUrl)
	if err != nil {
		return goods, err
	}

	if err := json.Unmarshal([]byte(imageUrl), &goods.ImagesUrls); err != nil {
		log.Fatal(err)
	}

	return goods, nil
}

func DeleteProduct(id int64) error {
	delForm, err := db.Prepare("DELETE FROM goods WHERE id=?")
	if err != nil {
		return err
	}
	defer delForm.Close()

	_, err = delForm.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func EditProduct(id int64, title string, description string, price float64, count int) error {
	stmt, err := db.Prepare("UPDATE goods SET title = ?, price = ?, goods_description = ?, count = ? WHERE id = ?;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, price, description, count, id)
	if err != nil {
		return err
	}
	return nil
}

func InsertGoodsIntoTable(title string, description string, price float64, count int, filenames []string) error {
	json_filenames, err := json.Marshal(filenames)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO goods (title, price, goods_description, count, json_urls) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, price, description, count, json_filenames)
	if err != nil {
		return err
	}

	return nil
}

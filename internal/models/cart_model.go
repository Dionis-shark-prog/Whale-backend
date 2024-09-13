package models

import (
	_ "database/sql"
	"encoding/json"
	"log"
)

type ProductInCart struct {
	Cart_ID    int64
	Client_ID  int64
	Product_ID int64
	Price      float64
	Title      string
	Count      int
	ImageUrl   []string
}

func InsertProductToCart(client_id int64, product_id int64, price float64) error {
	stmt, err := db.Prepare("INSERT INTO cart (client_id, product_id, price) VALUES (?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(client_id, product_id, price)
	if err != nil {
		return err
	}

	return nil
}

func GetAllFromCartWithID(client_id int64) ([]ProductInCart, error) {
	var cartProducts []ProductInCart

	stmt := "SELECT cart.cart_id, cart.client_id, cart.product_id, cart.price, goods.title, goods.count, goods.json_urls FROM cart JOIN goods ON cart.product_id = goods.id WHERE client_id=?"
	rows, err := db.Query(stmt, client_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productInCart ProductInCart
		var imageUrls string

		if err := rows.Scan(&productInCart.Cart_ID, &productInCart.Client_ID, &productInCart.Product_ID, &productInCart.Price, &productInCart.Title, &productInCart.Count, &imageUrls); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(imageUrls), &productInCart.ImageUrl); err != nil {
			log.Fatal(err)
		}

		cartProducts = append(cartProducts, productInCart)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cartProducts, nil
}

func DeleteProductFromCart(client_id int64, product_id int64) error {
	delForm, err := db.Prepare("DELETE FROM cart WHERE client_id=? AND product_id=?")
	if err != nil {
		return err
	}
	defer delForm.Close()

	_, err = delForm.Exec(client_id, product_id)
	if err != nil {
		return err
	}
	return nil
}

func ClearUserCart(client_id int64) error {
	delForm, err := db.Prepare("DELETE FROM cart WHERE client_id=?")
	if err != nil {
		return err
	}
	defer delForm.Close()

	_, err = delForm.Exec(client_id)
	if err != nil {
		return err
	}
	return nil
}

func IsProductByIDUserID(client_id int64, product_id int64) bool {
	stmt := "SELECT cart_id FROM cart WHERE client_id=? AND product_id=?"
	row := db.QueryRow(stmt, client_id, product_id)
	var dummy int64
	err := row.Scan(&dummy)
	return err == nil
}

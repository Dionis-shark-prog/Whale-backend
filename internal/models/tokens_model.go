package models

import "fmt"

func InsertAndCreateToken(id int64) error {
	stmt, err := db.Prepare("INSERT INTO clients_tokens (client_id, token_new) VALUES (?, UUID_TO_BIN(UUID()));")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func TakeTokenByID(id int64) (string, error) {
	var token string
	stmt := "SELECT token_new FROM clients_tokens WHERE client_id=?;"
	row := db.QueryRow(stmt, id)
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func TakeIDByToken(token string) (int64, error) {
	var id int64
	stmt := "SELECT client_id FROM clients_tokens WHERE token_new=?;"
	row := db.QueryRow(stmt, token)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func DeleteTokenByID(client_id int64) error {
	delForm, err := db.Prepare("DELETE FROM clients_tokens WHERE client_id=?")
	if err != nil {
		return err
	}
	defer delForm.Close()

	_, err = delForm.Exec(client_id)
	if err != nil {
		return err
	}

	fmt.Println("Deleted token by ID:", client_id)
	return nil
}

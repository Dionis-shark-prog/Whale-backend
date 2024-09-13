package models

import (
	"database/sql"
	"fmt"
)

type Clients struct {
	ID         int64
	Firstname  string
	Secondname string
	Email      string
	Number     string
}

func IsUserByName(name string, email string) (int, error) {
	// status: 0 - error but not "no rows"
	// status: 1 - no user found (we can insert current data)
	// status: 2 - user found
	stmt := "SELECT client_id FROM clients WHERE firstname=? AND email=?;"
	row := db.QueryRow(stmt, name, email)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		return 0, err
	} else if err == sql.ErrNoRows {
		return 1, nil
	} else {
		return 2, nil
	}
}

func TakeUserByToken(token string) (int64, error) {
	id, err := TakeIDByToken(token)
	if err != nil {
		return -1, err
	}

	stmt := "SELECT client_id FROM clients WHERE client_id=?;"
	row := db.QueryRow(stmt, id)
	var realID int64
	err = row.Scan(&realID)
	if err != nil {
		return -1, err
	}
	return realID, nil
}

func InsertUserToClients(firstname string, secondname *string, email string, hash []byte, number *string) error {
	stmt, err := db.Prepare("INSERT INTO clients (firstname, secondname, email, client_password, phone_number, user_role) VALUES (?, ?, ?, ?, ?, 'user');")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(firstname, secondname, email, hash, number)
	if err != nil {
		return err
	}

	return nil
}

func TakeHashByNameEmail(name string, email string) (string, error) {
	var hash string
	stmt := "SELECT client_password FROM clients WHERE firstname=? AND email=?;"
	row := db.QueryRow(stmt, name, email)
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func TakeIDByNameEmail(name string, email string) (int64, error) {
	var id int64
	stmt := "SELECT client_id FROM clients WHERE firstname=? AND email=?;"
	row := db.QueryRow(stmt, name, email)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func TakeUserDataByID(id int64) (Clients, error) {
	var userProfile Clients

	var sName, pNumber any

	stmt := "SELECT client_id, firstname, secondname, email, phone_number FROM clients WHERE client_id=?;"
	row := db.QueryRow(stmt, id)
	err := row.Scan(&userProfile.ID, &userProfile.Firstname, &sName, &userProfile.Email, &pNumber)
	if err != nil {
		return userProfile, err
	}

	if sName == nil {
		userProfile.Secondname = ""
	}
	if pNumber == nil {
		userProfile.Number = ""
	}

	return userProfile, nil
}

func TakeUserRole(name string, email string) (string, error) {
	var role string
	stmt := "SELECT user_role FROM clients WHERE firstname=? AND email=?;"
	row := db.QueryRow(stmt, name, email)
	err := row.Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func DeleteClientFromBase(client_id int64) error {
	delForm, err := db.Prepare("DELETE FROM clients WHERE client_id=?")
	if err != nil {
		return err
	}
	defer delForm.Close()

	_, err = delForm.Exec(client_id)
	if err != nil {
		return err
	}

	err = DeleteTokenByID(client_id)
	if err != nil {
		return err
	}

	fmt.Println("Deleted", client_id)
	return nil
}

func TakeTokenByNameEmail(name string, email string) (string, error) {
	id, err := TakeIDByNameEmail(name, email)
	if err != nil {
		return "", err
	}

	token, err := TakeTokenByID(id)
	if err != nil {
		return "", err
	}
	return token, nil
}

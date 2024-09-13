package models

import "database/sql"

type Admin struct {
	Name string
	Hash string
}

func IsAdminByID(id int64) error {
	stmt := "SELECT admin_id FROM admin WHERE admin_id=?;"
	row := db.QueryRow(stmt, id)
	var dummy int64
	err := row.Scan(&dummy)
	return err
}

func IsAdminByName(name string) (int, error) {
	stmt := "SELECT admin_id FROM admin WHERE admin_name=?;"
	row := db.QueryRow(stmt, name)
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

func InsertAdmin(name string, hash []byte) error {
	stmt, err := db.Prepare("INSERT INTO admin (admin_name, admin_password) VALUES (?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, hash)
	if err != nil {
		return err
	}

	return nil
}

func TakeIDByAdminName(name string) (int64, error) {
	var id int64
	stmt := "SELECT admin_id FROM admin WHERE admin_name=?;"
	row := db.QueryRow(stmt, name)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func TakeHashByAdminName(name string) (string, error) {
	var hash string
	stmt := "SELECT admin_password FROM admin WHERE admin_name=?;"
	row := db.QueryRow(stmt, name)
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

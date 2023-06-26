package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id,name,quantity,price from products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name,quantity,price FROM products where id=%v", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v',quantity=%v,price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such rows exists to")
	}
	return err
}
func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=%v", p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such rows exists")
	}
	return err
}
func (p *product) partialUpdate(db *sql.DB, param string) error {
	var query string
	if param == "Name" {
		query = fmt.Sprintf("update products set name='%v' where id=%v", p.Name, p.ID)
	}
	if param == "Quantity" {
		query = fmt.Sprintf("update products set quantity=%v where id=%v", p.Quantity, p.ID)
	}
	if param == "Price" {
		query = fmt.Sprintf("update products set price=%v where id=%v", p.Price, p.ID)
	}
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such rows exists/updating with the same value")
	}
	return err
}

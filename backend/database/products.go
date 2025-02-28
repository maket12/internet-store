package database

import (
	"github.com/google/uuid"
	"shop_backend/models"
)

func DBInitProducts() {
	gDB.Exec(`
  CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    image TEXT,
    price REAL,
    available INTEGER
  )
  `)
}

func GetAllProducts() ([]models.Product, error) {
	rows, err := gDB.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.Id, &p.Name, &p.Description, &p.Image, &p.Price, &p.Available)
		products = append(products, p)
	}
	return products, nil
}

func AddProduct(p models.Product) (string, error) {
	p.Id = uuid.New().String()
	_, err := gDB.Exec("INSERT INTO products (id, name, description, image, price, available) VALUES (?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.Description, p.Image, p.Price, p.Available)
	if err != nil {
		return "", err
	}
	return p.Id, nil
}

func UpdateProduct(p models.Product) error {
	_, err := gDB.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, available = ? WHERE id = ?", p.Name, p.Description, p.Image, p.Price, p.Available, p.Id)
	return err
}

func RemoveProduct(id string) error {
	_, err := gDB.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

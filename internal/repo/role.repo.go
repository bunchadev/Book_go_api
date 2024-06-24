package repo

import "Book_market_api/internal/database"

type RoleRepo struct{}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}

func GetRoleId(name string) (string, error) {
	query := `SELECT id FROM roles
              WHERE name = $1  
   `
	var id string
	err := database.DB().QueryRow(query, name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

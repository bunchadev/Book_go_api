package repo

import (
	"Book_market_api/internal/database"
	"Book_market_api/internal/models"
	"Book_market_api/utils"
	"errors"
	"fmt"
	"strings"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) GetUserPagination(page int, limit int, order string, field string, username string, email string) ([]models.UserPagination, error) {
	query := `
	   SELECT u.id,u.username,u.email,u.password,u.balance,u.login_enabled,u.depot_limit,u.created_at,r.name
	   FROM users u
	   INNER JOIN roles r ON r.id = u.role_id
	`
	var conditions []string
	var params []interface{}
	paramIndex := 1
	if username != "" {
		conditions = append(conditions, fmt.Sprintf("u.username ILIKE '%%' || $%d || '%%'", paramIndex))
		params = append(params, username)
		paramIndex++
	}
	if email != "" {
		conditions = append(conditions, fmt.Sprintf("u.email ILIKE '%%' || $%d || '%%'", paramIndex))
		params = append(params, email)
		paramIndex++
	}
	if len(conditions) > 0 {
		query += "WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", field, order, paramIndex, paramIndex+1)
	params = append(params, limit, (page-1)*limit)
	rows, err := database.DB().Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listUser []models.UserPagination
	for rows.Next() {
		var user models.UserPagination
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Balance, &user.Login_enabled, &user.Depot_limit, &user.Create_at, &user.Role)
		if err != nil {
			return nil, err
		}
		listUser = append(listUser, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listUser, nil
}

func (ur *UserRepo) CheckUserName(username string) bool {
	query := `SELECT id 
	          FROM users 
	          WHERE username = $1
			`
	id := ""
	err := database.DB().QueryRow(query, username).Scan(&id)
	if err != nil && id != "" {
		return true
	}
	if err == nil && id != "" {
		return true
	}
	return false
}

func (ur *UserRepo) CreateUser(user *models.UserCreate) error {
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	role_id, err := GetRoleId(user.Role)
	if err != nil {
		return err
	}
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username,email,password,role_id) 
	          VALUES ($1,$2,$3,$4)
			`
	_, err = tx.Exec(query, user.Username, user.Email, hashPassword, role_id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) UpdateUser(user *models.UpdateUser) error {
	query := `UPDATE users 
	          SET username = $1, email = $2 ,balance = $3,login_enabled = $4,depot_limit = $5,role_id = $6
	          WHERE id = $7
		     `
	role_id, err := GetRoleId(user.Role)
	if err != nil {
		return err
	}
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(query, user.Username, user.Email, user.Balance, user.Login_enabled, user.Depot_limit, role_id, user.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) DeleteUser(id string) error {
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM users
	          WHERE id = $1          
	         `
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) LoginUser(user *models.LoginUser) (*models.TokenResponse, error) {
	query := `SELECT u.id ,u.username,u.email,u.password,u.balance,u.depot_limit,r.name
              FROM users u
			  INNER JOIN roles r ON r.id = u.role_id
			  WHERE username = $1  	
	`
	var password string
	var userRes models.UserResponse
	err := database.DB().QueryRow(query, user.Username).Scan(
		&userRes.Id,
		&userRes.Username,
		&userRes.Email,
		&password,
		&userRes.Balance,
		&userRes.Depot_limit,
		&userRes.Role,
	)
	if checked := utils.CheckPasswordHash(user.Password, password); !checked {
		return nil, errors.New("password invalid")
	}
	if err != nil {
		return nil, err
	}
	access_token, err := utils.GenerateToken(userRes.Id, userRes.Role, 3)
	if err != nil {
		return nil, err
	}
	refresh_token, err := utils.GenerateToken(userRes.Id, userRes.Role, 5)
	if err != nil {
		return nil, err
	}
	tokenRes := models.TokenResponse{
		Access_token:  access_token,
		Refresh_token: refresh_token,
		Expires_in:    3600 * 3,
		User:          userRes,
	}
	return &tokenRes, nil
}

func (ur *UserRepo) GetNewToken(id string) (*models.TokenResponse, error) {
	query := `SELECT u.id,u.username,u.email,u.balance,u.depot_limit,r.name
              FROM users u
			  INNER JOIN roles r ON u.role_id = r.id
			  WHERE u.id = $1
	`
	var userRes models.UserResponse
	err := database.DB().QueryRow(query, id).Scan(
		&userRes.Id,
		&userRes.Username,
		&userRes.Email,
		&userRes.Balance,
		&userRes.Depot_limit,
		&userRes.Role,
	)
	if err != nil {
		return nil, err
	}
	access_token, err := utils.GenerateToken(userRes.Id, userRes.Role, 3)
	if err != nil {
		return nil, err
	}
	refresh_token, err := utils.GenerateToken(userRes.Id, userRes.Role, 5)
	if err != nil {
		return nil, err
	}
	tokenRes := models.TokenResponse{
		Access_token:  access_token,
		Refresh_token: refresh_token,
		Expires_in:    3600 * 3,
		User:          userRes,
	}
	return &tokenRes, nil
}

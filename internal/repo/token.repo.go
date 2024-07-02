package repo

import (
	"Book_market_api/internal/database"
	"time"
)

type TokenRepo struct{}

func NewTokeRepo() *TokenRepo {
	return &TokenRepo{}
}

func CheckToken(userId string) bool {
	query := `SELECT COUNT(*) FROM tokens
              WHERE user_id = $1	         
	`
	var count int
	err := database.DB().QueryRow(query, userId).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func DeleteToken(userId string) error {
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM tokens
	          WHERE user_id = $1
	`
	_, err = tx.Exec(query, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func CreateToken(access_token, refresh_token, userId string) error {
	if check := CheckToken(userId); check {
		DeleteToken(userId)
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
	query := `INSERT INTO tokens (user_id,access_token,refresh_token,access_token_expiry,refresh_token_expiry)
	          VALUES ($1,$2,$3,$4,$5)          
	`
	_, err = tx.Exec(
		query,
		userId,
		access_token,
		refresh_token,
		time.Now().Add(time.Minute*30),
		time.Now().Add(time.Hour*2),
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateToken(access_token, refresh_token, userId string) error {
	query := `UPDATE tokens
	          SET access_token = $1,refresh_token = $2,access_token_expiry = $3,updated_at = $4
			  WHERE user_id = $5            
	`
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(
		query,
		access_token,
		refresh_token,
		time.Now().Add(time.Minute*30),
		time.Now(),
		userId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tr *TokenRepo) TokenRetrieval(userId string) error {
	query := `UPDATE tokens
	          SET revoked = $1
			  WHERE user_id = $2
	`
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(query, true, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tr *TokenRepo) DeleteToken_v2(userId string) error {
	tx, err := database.DB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	query := `DELETE FROM tokens
	          WHERE user_id = $1
	`
	_, err = tx.Exec(query, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func CheckAccess_token(access_token string) bool {
	query := `SELECT access_token_expiry,revoked 
	          FROM tokens
			  WHERE access_token = $1
	`
	var checkTime time.Time
	checkRevoke := false
	err := database.DB().QueryRow(query, access_token).Scan(&checkTime, &checkRevoke)
	if err != nil {
		return false
	}
	if checkTime.Before(time.Now()) || checkRevoke {
		return false
	}
	return true
}

func CheckRefresh_token(refresh_token string) bool {
	query := `SELECT refresh_token_expiry,revoked 
	          FROM tokens
			  WHERE refresh_token = $1
	`
	var checkTime time.Time
	checkRevoke := false
	err := database.DB().QueryRow(query, refresh_token).Scan(&checkTime, &checkRevoke)
	if err != nil {
		return false
	}
	if checkTime.Before(time.Now()) || checkRevoke {
		return false
	}
	return true
}

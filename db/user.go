package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
	"net/http"
	customError	"github.com/nutcas3/my-ticko/custom_error"

	"github.com/nutcas3/my-ticko/db/model"
)

type DBUserInterface interface {
	CreateUser(username string, role model.Role) (int, error)
	GetUserById(id int) (*model.UserWithRoleList, error)
	GetUserByName(name string) (*model.UserWithRoleList, error)
}

func (pgdb *PostgresqlDB) CreateUser(username string, role model.Role) (int, error) {
	var userID int
	tx, err := pgdb.DB.Begin(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to make a transaction")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(context.Background())
		}
	}()
	err = tx.QueryRow(context.Background(), `
		INSERT INTO users (
			"username"
		)
		VALUES ($1)
		RETURNING id
	`,
		username,
	).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return 0, errors.New("Duplicate username")
			}
		}
		return 0, errors.Wrap(err, "Unable to create user")
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO user_roles(user_id,role_id) values ($1,$2)`, userID, role)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to add user role on create")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "Unable to commit a transaction")
	}
	return userID, nil
}
package db

import (
	"auth-service/transport/transport_http/handlers"
	"context"

	"go.uber.org/zap"
)

func (db *DB) Register(login string, pass_hash []byte) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		db.logger.Error("Failed to acquire connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `INSERT INTO "users"(login, pass_hash) VALUES ($1, $2)`, login, pass_hash)
	if err != nil {
		db.logger.Error("Failed to register user", zap.Error(err))
		return err
	}

	return nil
}

func (db *DB) GetUser(login string) (handlers.User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		db.logger.Error("Failed to acquire connection", zap.Error(err))
		return handlers.User{}, err
	}
	defer conn.Release()

	var user handlers.User

	err = conn.QueryRow(context.Background(), `SELECT * FROM "users" WHERE login = $1`, login).Scan(&user.ID, &user.Login, &user.Role, &user.Hash)
	if err != nil {
		db.logger.Error("Failed to find user", zap.Error(err))
		return handlers.User{}, err
	}

	return user, nil
}

func (db *DB) ChangePassword(login string, pass_hash string, new_hash []byte) error {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		db.logger.Error("Failed to acquire connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `UPDATE "users" SET pass_hash = $1 WHERE login = $2 AND pass_hash = $3`, new_hash, login, pass_hash)
	if err != nil {
		db.logger.Error("Failed to register user", zap.Error(err))
		return err
	}

	return nil
}

package repository

import (
	"context"
	"fmt"
)

type User struct {
	Id             int    `json:"id" db:"id"`
	Login          string `json:"login" db:"login"`
	HashedPassword string `json:"hashed_password" db:"hashed_password"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
}

// Login ищет пользователя по логину и паролю в БД
func (r *Repository) Login(ctx context.Context, login, hashedPassword string) (u User, err error) {
	// запрос к БД
	row := r.pool.QueryRow(ctx, "select id, login, name, surname from users where login = $1"+
		"and hashed_password = $2", login, hashedPassword)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	// записываем в структуру
	err = row.Scan(&u.Id, &u.Login, &u.Name, &u.Surname)
	if err != nil {
		err = fmt.Errorf("failed to query data: %w", err)
		return
	}
	return
}

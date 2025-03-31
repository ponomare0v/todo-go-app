package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ponomare0v/todo-go-app/pkg/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	//Метод QueryRaw возвращает объект raw, то есть он хранит в себе информацию о возвращаемой строке из базы (в нашем случае запрос в
	//  бд возвращает одну строку со значением поля id).
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id from %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

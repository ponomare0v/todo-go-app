package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ponomare0v/todo-go-app/pkg/models"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

//
//
//
//

func (r *TodoListPostgres) Create(userId int, list models.TodoList) (int, error) {
	tx, err := r.db.Begin() //Для создании транзакции в объектах БД есть метод Begin
	if err != nil {
		return 0, err
	}

	// Запрос для создании записи в таблице todo_lists, возвращая id нового списка.
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Вставка в таблицу  users_lists, в которой свяжем id пользователя и id нового списка.
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id) //Для простого выполнения запроса, без чтения возвращаемой инфоормации - метод Exec.
	if err != nil {
		tx.Rollback() //В случае ошибок - вызываем метод Rollback у транзакции, который откатывает все изменения базы данных до начала выполнения транзакции.
		return 0, err
	}

	return id, tx.Commit() // После выполнения транзакции вызовем метод Commit, который применит наши изменения к БД и закончит транзакцию.
}

//
//
//
//
//

func (r *TodoListPostgres) GetAll(userId int) ([]models.TodoList, error) {
	var lists []models.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId) //db.Select() - работает аналогично с методом db.Get() только применяется при выборке больше одного элемента
	//  и для записи в слайс. Нужно добавить теги db в наши модели, чтобы иметь возможность сделать выборки из базы

	return lists, err
}

//
//
//
//
//

func (r *TodoListPostgres) GetById(userId, listId int) (models.TodoList, error) {
	var list models.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2", todoListsTable, usersListsTable) //добавили доп условие для проверки id листа
	err := r.db.Get(&list, query, userId, listId)                                                                                                                                               // метод get

	return list, err
}

//
//
//
//
//

func (r *TodoListPostgres) Delete(userId, listId int) error {

	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input models.UpdateListInput) error {
	//Инициализируем три переменных - 1 слайс строк 2 слайс интерфейсов 3 id аргумента
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	//  Далее будем выполнять проверку полей, если они не нил, значит будем добавлять в слайсы
	// элементы для формирования запроса в базу с их обновлением. В слайс строк будем записывать присвоение
	//  полю title а после знака = будем записывать значение для placeholder. В слайс значения добавим
	//  само значение поля titile.
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// setQuery, в которой соединяем элементы слайса строк в одну строку
	setQuery := strings.Join(setValues, ", ")

	// В качестве аргументов функции sprintf передим названия таблиц, строку со значениями на обновление полей
	// а также значение ardId и argId + 1 для placeholder. Вот как будет выглядеть строка для setquery для всех трех варинатов:
	//title=$1
	//description=$1
	//title=$1, description=$2

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	//В слайс аргументов добавим еще два элемента id пользователя и списка, а также залогируем запрос и аргументы в консоль
	// для наглядности.
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

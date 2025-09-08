package db

import (
	"database/sql"
	"fmt"
)

// Task - структура для парсинга таблица - задача из БД.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask - функция для добавления задачи в таблицу БД.
// Принимает структуру Task и возвращает id добавленной задачи и ошибку.
func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// TasksWithSearch - функция для поиск задачи по определённым параметрам.
// Принимает лимит количества возвращаемых строк таблицы, строку поиска
// и булевую переменную:
// true - если в строку поиска передана дата (тогда поиск будет осуществляться по дате задачи),
// false - если в строку поиска передана подстрака, по которой будет осуществляться поиск задачи в парметрах  - название и комментарий задачи.
// Возвращает список задач отсортированных по дате и ошибку.
func TasksWithSearch(limit int, search string, date bool) ([]*Task, error) {
	var query string
	if date {
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :search  ORDER BY date ASC LIMIT :limit"

	} else {
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC LIMIT :limit"
		search = fmt.Sprintf("%%%s%%", search)
	}
	rows, err := db.Query(query,
		sql.Named("search", search),
		sql.Named("limit", limit))

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return Tasks(limit, rows)

}

// TasksWithOutSearch - функция, которая принимает лимит количества возвращаемых строк
// и возвращает список задач отсортированнх по дате от меньшего к большиму и ошибку.
func TasksWithOutSearch(limit int) ([]*Task, error) {
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return Tasks(limit, rows)
}

// Tasks - функция для поиска задач.
// Принимает лимит количества задач и строики из таблицы БД,
// возвращает список задач и таблицу.
// Вызывается в функциях TasksWithOutSearch и TasksWithSearch данного пакета.
func Tasks(limit int, rows *sql.Rows) ([]*Task, error) {
	tasks := make([]*Task, 0, limit)

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)

	}

	return tasks, nil
}

// GetTask - функция для поиска задачи по её id.
// Принимает строку - id задачи,
// возвращает задачу и ошибку
func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask - функция редактирования задачи.
// Принимает ссылку на структуру  Task с изменёнными параметрами, возвращает ошибку.
func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

// UpdateDate - функция для редактирования даты задачи.
// Принимает строку - дата, на которую необходимо заменить текущее значение и id задачи,
// возвращает ошибку.
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("id", id),
		sql.Named("date", next))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

// DeleteTask  - функция, которая удаляет задачу.
// Принмает id задачи, котроую необходимо удалить,
// возвращает ошибку
func DeleteTask(id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`not found`)
	}
	return nil
}

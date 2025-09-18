// Пакет db реализует работу с базой данных планировщика задач.
// База данных реализуется на SQLite.
// Этот пакет включает в себя такие функции, как:
// 1)Функция создания базы данных, если её ещё нет, и подключение к ней;
// 2)Функция закрытия подключение к базе данных;
// 3)Функция добавления задачи в таблицу БД;
// 4)Функции поиска задач в БД;
// 5)Функция поиска задачи по id;
// 6)Функции редактирования задачи;
// 7)Функция удаления задачи из БД.
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// schema - строка заздания таблицы БД.
const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);

	CREATE INDEX scheduler_date ON scheduler (date);`

// db - интерфейс для работы с БД.
var db *sql.DB

// Init - функция инициализации БД.
// Создаёт базу данных, если она ещё не создана
// и реализует подключение к базе данных.
// Принимает строку - название файла с расширением .db,
// возвращает ошибку.
func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil

}

// Close - функция для закрытия базы данных вне пакета.
func Close() {
	db.Close()
}

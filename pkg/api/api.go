// Пакет api реализует API планировщика задач.
// Пакет реализует хендлеры и их инициализацию.
// Реализованы следующие обработчики запросов:
// 1)"/api/signin" - реализует метод Post, предназначен для аутентификации пользователей.
// 2)"/api/tasks" - реализует метод Get, возвращает клиенту json - список задач или ошибку.
// 3)"/api/task" - реализует методы Get, Post, Put, Delete.
//
//			метод Get - возвращает клиенту задачу по её id.
//			метод Post - добавляет новую задачу в БД.
//	     	метод Put - редактирует задачу.
//	     	метод Delete - удаляет задачу из БД.
//
// 4)"/api/nextdate" - релизует метод Get, изменят дату задачи, если установлено правило повторения.
// 5)"/api/task/done" - реализует метод Post, предназначена для перевода задачи в статус - выполнена,
// если для задачи указано правило повторения - переопределяет дату задачи на следующую, если правило не указано - удаляет задачу.
package api

import (
	"encoding/json"
	"net/http"
)

// DataFormat устанавливает форматы
const DataFormat = "20060102"

// Init - инициализация обработчиков, вызывается в функции main.
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/signin", authHandler)

}

// taskHandler - устанавливает маршрут запросов для "/api/task".
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}

}

// errorResp - структура для формирования JSON, котороый сообщает клиенту об ошибке.
type errorResp struct {
	Error string `json:"error"`
}

// writeJson - метод ля отправки ответа клиаету в формате JSON.
// Принимает http.ResponseWriter, данные для передачи и булевую пременную:
// true - формируется ответ со статусом 200,
// false - формирут ответ со статусом 400.
func writeJson(w http.ResponseWriter, data any, ok bool) {

	if !ok {
		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

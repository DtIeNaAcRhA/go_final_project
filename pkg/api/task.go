package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DtIeNaAcRhA/go_final_project/pkg/db"
)

// getTaskHandler - хендлер который в теле ответа возвращает задачу, найденную по id, который передаётся в параметре URL запроса.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var err error
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, errorResp{Error: "Не указан идентификатор"}, false)
		return
	}
	task, err = db.GetTask(id)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}
	writeJson(w, task, true)
}

// updateTaskHandler - хендлер для редактирования задачи.
// В теле запроса ожидает JSON с параметрами задачи,
// в теле ответа возращает пустую структуру при успешном редактировании задачи или ошибку - error.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var task *db.Task

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, errorResp{Error: "Ошибка десериализации JSON"}, false)
		return
	}

	if task.ID == "" {
		writeJson(w, errorResp{Error: "Не указан идентификатор"}, false)
		return
	}

	if task.Title == "" {
		writeJson(w, errorResp{Error: "Не указан заголовок задачи"}, false)
		return
	}

	err = checkDate(task)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}
	writeJson(w, make(map[string]interface{}), true)

}

// deleteTaskHandler - хендлер для удаления задачи по id, который передаётся в параметре URL запроса.
// В теле ответа возращает пустую структуру при успешном удалении задачи или ошибку - error.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, errorResp{Error: "Не указан идентификатор"}, false)
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}
	writeJson(w, make(map[string]interface{}), true)
}

// taskDoneHandler -  хендлер для перевода задачи в статус - выполнена.
// Реализует метод Post.
// Если для задачи указано правило повторения - переопределяет дату задачи на следующую,
// если правило не указано - удаляет задачу.
// В теле ответа возращает пустую структуру при успешном исполнении или ошибку - error.
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	var err error

	if r.Method != http.MethodPost {
		writeJson(w, errorResp{Error: "Method not allowed"}, false)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, errorResp{Error: "Не указан идентификатор"}, false)
		return
	}
	task, err = db.GetTask(id)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, errorResp{Error: err.Error()}, false)
			return
		}
		writeJson(w, make(map[string]interface{}), true)
		return
	}
	now := time.Now()
	result, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}

	err = db.UpdateDate(result, task.ID)
	if err != nil {
		writeJson(w, errorResp{Error: err.Error()}, false)
		return
	}
	writeJson(w, make(map[string]interface{}), true)
}
